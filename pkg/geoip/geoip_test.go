package geoip

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"TOomaAh/emby_exporter_go/pkg/logger"
)

func TestNoGeoIPLookup(t *testing.T) {
	record, err := NewNoGeoIP().Lookup("8.8.8.8")
	if !errors.Is(err, ErrGeoIPDisabled) {
		t.Errorf("err = %v, want %v", err, ErrGeoIPDisabled)
	}
	if record != nil {
		t.Errorf("record = %+v, want nil", record)
	}
}

func TestGeoIPManagerLookupWithoutDatabase(t *testing.T) {
	g := &GeoIPManager{logger: logger.New("error")}

	tests := []struct {
		ip        string
		wantLocal bool
		wantErr   error
	}{
		{"not an ip", false, ErrInvalidIPAddress},
		{"127.0.0.1", true, nil},
		{"192.168.1.20", true, nil},
		{"10.0.0.1", true, nil},
		// A public IP needs the database: no panic when it is not loaded (or closed).
		{"8.8.8.8", false, ErrDatabaseNotInitialized},
	}

	for _, tt := range tests {
		record, err := g.Lookup(tt.ip)
		if tt.wantErr != nil {
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Lookup(%q) err = %v, want %v", tt.ip, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("Lookup(%q) err = %v", tt.ip, err)
			continue
		}
		if tt.wantLocal && record.City != "Local" {
			t.Errorf("Lookup(%q) = %+v, want the local record", tt.ip, record)
		}
	}
}

func TestGeoIPManagerCloseTwice(t *testing.T) {
	g := &GeoIPManager{logger: logger.New("error")}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.mmdb")
	dst := filepath.Join(dir, "dst.mmdb")

	data := bytes.Repeat([]byte("geoip"), 1<<18)
	if err := os.WriteFile(src, data, 0o644); err != nil {
		t.Fatal(err)
	}
	// dst already exists and is longer: it must be fully replaced.
	if err := os.WriteFile(dst, bytes.Repeat([]byte("old"), 1<<19), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("dst has %d bytes, want the %d bytes of src", len(got), len(data))
	}
}

func TestCopyFileMissingSource(t *testing.T) {
	dir := t.TempDir()
	if err := copyFile(filepath.Join(dir, "missing"), filepath.Join(dir, "dst")); err == nil {
		t.Error("want an error for a missing source")
	}
}
