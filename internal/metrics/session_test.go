package metrics

import (
	"net/http"
	"testing"

	"TOomaAh/emby_exporter_go/pkg/geoip"
	"TOomaAh/emby_exporter_go/pkg/logger"
)

const sessionsJSON = `[
	{
		"UserName": "alice", "Client": "Web", "RemoteEndPoint": "1.2.3.4",
		"PlayState": {"PlayMethod": "DirectPlay", "PositionTicks": 18000000000},
		"NowPlayingItem": {"Name": "Pilot", "Type": "Episode", "SeriesName": "Breaking Bad", "SeasonName": "Season 1", "RunTimeTicks": 36000000000, "Bitrate": 2048000}
	},
	{
		"UserName": "bob", "Client": "TV", "RemoteEndPoint": "5.6.7.8",
		"PlayState": {"PlayMethod": "Transcode"}
	},
	{
		"UserName": "idle", "Client": "Phone"
	}
]`

func sessionsMux(body string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/Sessions", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	})
	return mux
}

// Regression: with GeoIP disabled (the default) and a session without
// NowPlayingItem, Collect used to panic on nil pointers.
func TestSessionCollectorWithGeoIPDisabled(t *testing.T) {
	c := NewSessionCollector(newFakeEmby(t, sessionsMux(sessionsJSON)), geoip.NewNoGeoIP(), logger.New("error"))
	got := gather(t, c)

	if count := got["emby_sessions_count"]; len(count) != 1 || count[0].value != 2 {
		t.Fatalf("emby_sessions_count = %+v, want 2 (the idle session is not playing)", count)
	}

	sessions := map[string]map[string]string{}
	for _, s := range got["emby_sessions"] {
		sessions[s.labels["username"]] = s.labels
	}
	if len(sessions) != 2 {
		t.Fatalf("got %d emby_sessions, want 2: %+v", len(sessions), got["emby_sessions"])
	}

	alice := sessions["alice"]
	for label, want := range map[string]string{
		"nowPlayingItemName": "Pilot",
		"nowPlayingItemType": "Episode",
		"tvshow":             "Breaking Bad",
		"season":             "Season 1",
		"percentPlayback":    "50",
		"playMethod":         "Direct Play",
		"mediaDuration":      "01h00m00s",
		"currentPlayTime":    "00h30m00s",
		"city":               "N/A",
		"countryCode":        "N/A",
	} {
		if alice[label] != want {
			t.Errorf("alice %s = %q, want %q", label, alice[label], want)
		}
	}

	bob := sessions["bob"]
	for label, want := range map[string]string{
		"nowPlayingItemName": "N/A",
		"nowPlayingItemType": "N/A",
		"tvshow":             "",
		"percentPlayback":    "0",
		"city":               "N/A",
	} {
		if bob[label] != want {
			t.Errorf("bob %s = %q, want %q", label, bob[label], want)
		}
	}

	if len(got["emby_sessions_bitrate"]) != 2 {
		t.Errorf("got %d emby_sessions_bitrate, want 2", len(got["emby_sessions_bitrate"]))
	}
}

func TestSessionCollectorNoSession(t *testing.T) {
	c := NewSessionCollector(newFakeEmby(t, sessionsMux(`[]`)), geoip.NewNoGeoIP(), logger.New("error"))
	got := gather(t, c)

	if count := got["emby_sessions_count"]; len(count) != 1 || count[0].value != 0 {
		t.Errorf("emby_sessions_count = %+v, want 0", count)
	}
	if len(got["emby_sessions"]) != 0 {
		t.Errorf("want no emby_sessions, got %+v", got["emby_sessions"])
	}
}

func TestSessionCollectorEmbyError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/Sessions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	c := NewSessionCollector(newFakeEmby(t, mux), geoip.NewNoGeoIP(), logger.New("error"))

	if got := gather(t, c); len(got) != 0 {
		t.Errorf("want no metric when Emby fails, got %+v", got)
	}
}
