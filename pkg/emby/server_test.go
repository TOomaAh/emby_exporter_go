package emby

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"TOomaAh/emby_exporter_go/pkg/logger"
)

// newTestServer returns a Server talking to a fake Emby that answers with handler.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*Server, *httptest.Server) {
	t.Helper()
	fake := httptest.NewServer(handler)
	t.Cleanup(fake.Close)

	host, port, err := net.SplitHostPort(strings.TrimPrefix(fake.URL, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	s := NewServer(&ServerInfo{Hostname: "http://" + host, Port: port, Token: "token"}, logger.New("error"))
	return s, fake
}

func TestRequestsSendToken(t *testing.T) {
	s, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Emby-Token"); got != "token" {
			t.Errorf("X-Emby-Token = %q, want token", got)
		}
		w.Write([]byte(`{}`))
	})
	if _, err := s.GetServerInfo(); err != nil {
		t.Fatal(err)
	}
}

func TestPingUpdatesReachability(t *testing.T) {
	var status atomic.Int32
	status.Store(http.StatusOK)
	s, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(int(status.Load()))
	})

	if s.IsReachable() {
		t.Fatal("a new server must start unreachable")
	}
	if err := s.Ping(); err != nil || !s.IsReachable() {
		t.Fatalf("Ping() = %v, reachable = %v; want nil, true", err, s.IsReachable())
	}

	status.Store(http.StatusServiceUnavailable)
	if err := s.Ping(); err == nil || s.IsReachable() {
		t.Fatalf("Ping() = %v, reachable = %v; want an error, false", err, s.IsReachable())
	}
}

// Only a network error means the server is down: an HTTP or decoding error
// means it answered, so it must stay reachable (no ping loop flapping on a bad token).
func TestGetKeepsServerReachableOnHTTPErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"unauthorized", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) }},
		{"server error", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }},
		{"invalid JSON", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{not json`)) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := newTestServer(t, tt.handler)
			s.reachable.Store(true)

			if _, err := s.GetSessions(); err == nil {
				t.Fatal("want an error")
			}
			if !s.IsReachable() {
				t.Error("server marked unreachable after an HTTP error")
			}
		})
	}
}

func TestGetMarksServerUnreachableOnNetworkError(t *testing.T) {
	s, fake := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	s.reachable.Store(true)
	fake.Close()

	if _, err := s.GetSessions(); err == nil {
		t.Fatal("want an error")
	}
	if s.IsReachable() {
		t.Error("server still reachable after a network error")
	}
}

func TestGetSessionsDecodes(t *testing.T) {
	s, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Sessions" {
			t.Errorf("path = %s, want /Sessions", r.URL.Path)
		}
		w.Write([]byte(`[{"UserName":"alice","Client":"Web","PlayState":{"PlayMethod":"DirectPlay"}}]`))
	})

	sessions, err := s.GetSessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].UserName != "alice" || !sessions[0].HasPlayMethod() {
		t.Errorf("sessions = %+v", sessions)
	}
}

func TestGetLibrarySizeQuery(t *testing.T) {
	tests := []struct {
		contentType   string
		wantTypes     string
		wantNoFolders bool
	}{
		{"movies", "Movie", false},
		{"tvshows", "Series", false},
		{"music", "MusicArtist", false},
		{"playlists", "Playlist", false},
		{"homevideos", "Video,Photo", false},
		{"unknown", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			var (
				mu    sync.Mutex
				path  string
				query url.Values
			)
			s, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				path, query = r.URL.Path, r.URL.Query()
				mu.Unlock()
				w.Write([]byte(`{"Items":[],"TotalRecordCount":42}`))
			})

			size, err := s.GetLibrarySize("lib-1", tt.contentType)
			if err != nil {
				t.Fatal(err)
			}
			if size != 42 {
				t.Errorf("size = %d, want 42", size)
			}

			mu.Lock()
			defer mu.Unlock()
			if path != "/Items" {
				t.Errorf("path = %s, want /Items (no user scope)", path)
			}

			for key, want := range map[string]string{"ParentId": "lib-1", "Recursive": "true", "Limit": "0"} {
				if got := query.Get(key); got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
			if got := query.Get("IncludeItemTypes"); got != tt.wantTypes {
				t.Errorf("IncludeItemTypes = %q, want %q", got, tt.wantTypes)
			}
			if got := query.Get("IsFolder") == "false"; got != tt.wantNoFolders {
				t.Errorf("IsFolder=false sent = %v, want %v", got, tt.wantNoFolders)
			}
		})
	}
}

func TestStartPingLoopStopsWhenReachable(t *testing.T) {
	var calls atomic.Int32
	s, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// The first two pings fail.
		if calls.Add(1) <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	done := make(chan struct{})
	go func() {
		s.StartPingLoop(time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("ping loop did not stop once the server answered")
	}
	if !s.IsReachable() {
		t.Error("server not reachable after the ping loop")
	}
}

func TestStartPingLoopRunsOnce(t *testing.T) {
	s, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})

	// Simulate a loop already running: a second call must return immediately.
	s.pinging.Store(true)
	done := make(chan struct{})
	go func() {
		s.StartPingLoop(time.Hour)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("StartPingLoop started a second loop")
	}
}
