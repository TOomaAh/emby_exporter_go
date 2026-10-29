package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNewConfigDefaults(t *testing.T) {
	cfg, err := NewConfig(writeConfig(t, "server:\n  token: abc\n"))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Server.Hostname != "http://localhost" {
		t.Errorf("Hostname = %q, want http://localhost", cfg.Server.Hostname)
	}
	if cfg.Server.Port != "8096" {
		t.Errorf("Server.Port = %q, want 8096", cfg.Server.Port)
	}
	if cfg.Server.Token != "abc" {
		t.Errorf("Token = %q, want abc", cfg.Server.Token)
	}
	if cfg.Exporter.Port != 9210 {
		t.Errorf("Exporter.Port = %d, want 9210", cfg.Exporter.Port)
	}
	if cfg.Options.RetryInterval != 10 {
		t.Errorf("RetryInterval = %d, want 10", cfg.Options.RetryInterval)
	}
	if !cfg.Options.HealthCheck {
		t.Error("HealthCheck = false, want true by default")
	}
	if cfg.Options.GeoIP {
		t.Error("GeoIP = true, want false by default")
	}
}

func TestNewConfigOverrides(t *testing.T) {
	cfg, err := NewConfig(writeConfig(t, `
server:
  url: "https://emby.example.com"
  port: "443"
  token: abc
exporter:
  port: 9999
options:
  retryInterval: 30
  healthcheck: false
  geoip: true
`))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Server.Hostname != "https://emby.example.com" {
		t.Errorf("Hostname = %q, want https://emby.example.com", cfg.Server.Hostname)
	}
	if cfg.Server.Port != "443" {
		t.Errorf("Server.Port = %q, want 443", cfg.Server.Port)
	}
	if cfg.Exporter.Port != 9999 {
		t.Errorf("Exporter.Port = %d, want 9999", cfg.Exporter.Port)
	}
	if cfg.Options.RetryInterval != 30 {
		t.Errorf("RetryInterval = %d, want 30", cfg.Options.RetryInterval)
	}
	if cfg.Options.HealthCheck {
		t.Error("HealthCheck = true, want false")
	}
	if !cfg.Options.GeoIP {
		t.Error("GeoIP = false, want true")
	}
}

func TestNewConfigAddsHTTPScheme(t *testing.T) {
	cfg, err := NewConfig(writeConfig(t, "server:\n  url: 192.168.1.10\n  token: abc\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Hostname != "http://192.168.1.10" {
		t.Errorf("Hostname = %q, want http://192.168.1.10", cfg.Server.Hostname)
	}
}

// An empty section decodes as null and resets its fields: the values that must
// never be zero are restored, the others keep their zero value.
func TestNewConfigEmptySection(t *testing.T) {
	cfg, err := NewConfig(writeConfig(t, "server:\n  token: abc\nexporter:\noptions:\n"))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Options.RetryInterval != 10 {
		t.Errorf("RetryInterval = %d, want 10", cfg.Options.RetryInterval)
	}
	if cfg.Exporter.Port != 9210 {
		t.Errorf("Exporter.Port = %d, want 9210", cfg.Exporter.Port)
	}
}

func TestNewConfigInvalidValues(t *testing.T) {
	cfg, err := NewConfig(writeConfig(t, "server:\n  token: abc\nexporter:\n  port: -1\noptions:\n  retryInterval: -5\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Exporter.Port != 9210 {
		t.Errorf("Exporter.Port = %d, want 9210", cfg.Exporter.Port)
	}
	if cfg.Options.RetryInterval != 10 {
		t.Errorf("RetryInterval = %d, want 10", cfg.Options.RetryInterval)
	}
}

// The defaults must not leak between two loads.
func TestNewConfigDoesNotMutateDefaults(t *testing.T) {
	if _, err := NewConfig(writeConfig(t, "server:\n  token: abc\n  port: \"1234\"\n")); err != nil {
		t.Fatal(err)
	}
	cfg, err := NewConfig(writeConfig(t, "server:\n  token: abc\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "8096" {
		t.Errorf("Server.Port = %q, want 8096", cfg.Server.Port)
	}
}

func TestNewConfigErrors(t *testing.T) {
	if _, err := NewConfig(filepath.Join(t.TempDir(), "missing.yml")); err == nil {
		t.Error("missing file: want an error")
	}
	if _, err := NewConfig(writeConfig(t, "server: [not, a, map\n")); err == nil {
		t.Error("invalid YAML: want an error")
	}
}
