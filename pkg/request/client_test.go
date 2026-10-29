package request

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type testClient struct {
	baseURL *url.URL
}

func (c *testClient) GetBaseURL() *url.URL { return c.baseURL }

func (c *testClient) ApplyAuthentication(r *http.Request) error {
	r.Header.Set("X-Test-Token", "secret")
	return nil
}

func (c *testClient) SetHeaders(headers http.Header) {
	headers.Set("Accept", "application/json")
}

func (c *testClient) GetClient() *http.Client { return http.DefaultClient }

func newManager(t *testing.T, handler http.HandlerFunc) *RequestManager {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	return NewRequestManager(&testClient{baseURL: u})
}

func get(t *testing.T, rm *RequestManager, v interface{}) error {
	t.Helper()
	req, err := rm.NewRequest(http.MethodGet, "/path?a=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	return rm.Do(req, v)
}

func TestNewRequestAppliesHeadersAndAuthentication(t *testing.T) {
	rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/path" || r.URL.Query().Get("a") != "1" {
			t.Errorf("URL = %s, want /path?a=1", r.URL)
		}
		if got := r.Header.Get("X-Test-Token"); got != "secret" {
			t.Errorf("X-Test-Token = %q, want secret", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q, want application/json", got)
		}
	})
	if err := get(t, rm, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDoDecodesJSON(t *testing.T) {
	rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Name":"emby","Count":3,"Ignored":true}`))
	})

	var v struct {
		Name  string
		Count int
	}
	if err := get(t, rm, &v); err != nil {
		t.Fatal(err)
	}
	if v.Name != "emby" || v.Count != 3 {
		t.Errorf("decoded %+v, want {Name:emby Count:3}", v)
	}
}

func TestDoWithoutTarget(t *testing.T) {
	rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Emby Server"))
	})
	if err := get(t, rm, nil); err != nil {
		t.Fatalf("Do(nil) = %v, want nil", err)
	}
}

func TestDoInvalidJSON(t *testing.T) {
	rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{not json`))
	})
	var v map[string]any
	if err := get(t, rm, &v); err == nil {
		t.Fatal("want a decoding error")
	}
}

func TestDoStatusCodes(t *testing.T) {
	tests := []struct {
		status int
		want   error // nil means "any error"
	}{
		{http.StatusBadRequest, Error400BadRequest},
		{http.StatusUnauthorized, Error401Unauthorized},
		{http.StatusForbidden, Error403Forbidden},
		{http.StatusNotFound, Error404NotFound},
		{http.StatusInternalServerError, Error500Internal},
		{http.StatusBadGateway, nil},
		{http.StatusServiceUnavailable, nil},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(`{"error":"x"}`))
			})

			var v map[string]any
			err := get(t, rm, &v)
			if err == nil {
				t.Fatalf("status %d: want an error", tt.status)
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Errorf("status %d: err = %v, want %v", tt.status, err, tt.want)
			}
		})
	}
}

func TestDoNetworkError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	u, _ := url.Parse(srv.URL)
	srv.Close()

	rm := NewRequestManager(&testClient{baseURL: u})
	req, err := rm.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = rm.Do(req, nil)
	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatalf("err = %v, want a *url.Error", err)
	}
}

func TestDoFile(t *testing.T) {
	rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file content"))
	})
	req, err := rm.NewRequest(http.MethodGet, "/file", nil)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := rm.DoFile(req, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "file content" {
		t.Errorf("content = %q, want %q", out.String(), "file content")
	}
}

func TestDoFileStatusError(t *testing.T) {
	rm := newManager(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	req, err := rm.NewRequest(http.MethodGet, "/file", nil)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := rm.DoFile(req, &out); !errors.Is(err, Error401Unauthorized) {
		t.Errorf("err = %v, want %v", err, Error401Unauthorized)
	}
}
