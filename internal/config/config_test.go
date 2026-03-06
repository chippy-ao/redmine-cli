package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Fatal("DefaultConfigPath() returned empty string")
	}
	if !filepath.IsAbs(path) {
		t.Fatalf("DefaultConfigPath() returned non-absolute path: %s", path)
	}
	if filepath.Base(path) != "config.yaml" {
		t.Fatalf("DefaultConfigPath() should end with config.yaml, got: %s", path)
	}
}

func TestLoadConfig_NonexistentFile(t *testing.T) {
	cfg, err := LoadConfig("/tmp/redmine-cli-test-nonexistent/config.yaml")
	if err != nil {
		t.Fatalf("LoadConfig with nonexistent file should not error, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig should return non-nil Config")
	}
	if cfg.Profiles == nil {
		t.Fatal("Profiles map should be initialized (not nil)")
	}
	if len(cfg.Profiles) != 0 {
		t.Fatalf("Profiles should be empty, got %d entries", len(cfg.Profiles))
	}
}

func TestSaveAndLoadConfig_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.yaml")

	cfg := &Config{
		DefaultProfile: "work",
		Profiles: map[string]Profile{
			"work": {
				URL:    "https://redmine.company.com",
				APIKey: "abc123",
			},
			"oss": {
				URL:    "https://redmine.example.org",
				APIKey: "def456",
			},
		},
	}

	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Check directory permissions
	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("failed to stat directory: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0700 {
		t.Fatalf("directory permissions should be 0700, got %o", perm)
	}

	// Check file permissions
	finfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	if perm := finfo.Mode().Perm(); perm != 0600 {
		t.Fatalf("file permissions should be 0600, got %o", perm)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.DefaultProfile != cfg.DefaultProfile {
		t.Fatalf("DefaultProfile mismatch: got %q, want %q", loaded.DefaultProfile, cfg.DefaultProfile)
	}
	if len(loaded.Profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(loaded.Profiles))
	}

	work := loaded.Profiles["work"]
	if work.URL != "https://redmine.company.com" {
		t.Fatalf("work URL mismatch: got %q", work.URL)
	}
	if work.APIKey != "abc123" {
		t.Fatalf("work APIKey mismatch: got %q", work.APIKey)
	}
}

func TestGetProfile_Default(t *testing.T) {
	cfg := &Config{
		DefaultProfile: "work",
		Profiles: map[string]Profile{
			"work": {URL: "https://redmine.company.com", APIKey: "abc123"},
		},
	}

	p, err := cfg.GetProfile("")
	if err != nil {
		t.Fatalf("GetProfile with empty name should use default: %v", err)
	}
	if p.URL != "https://redmine.company.com" {
		t.Fatalf("URL mismatch: got %q", p.URL)
	}
}

func TestGetProfile_Named(t *testing.T) {
	cfg := &Config{
		DefaultProfile: "work",
		Profiles: map[string]Profile{
			"work": {URL: "https://redmine.company.com", APIKey: "abc123"},
			"oss":  {URL: "https://redmine.example.org", APIKey: "def456"},
		},
	}

	p, err := cfg.GetProfile("oss")
	if err != nil {
		t.Fatalf("GetProfile with valid name should not error: %v", err)
	}
	if p.URL != "https://redmine.example.org" {
		t.Fatalf("URL mismatch: got %q", p.URL)
	}
}

func TestGetProfile_EmptyConfig(t *testing.T) {
	cfg := &Config{
		Profiles: map[string]Profile{},
	}

	_, err := cfg.GetProfile("")
	if err == nil {
		t.Fatal("GetProfile with empty config and no name should return error")
	}
}

func TestGetProfile_UnknownName(t *testing.T) {
	cfg := &Config{
		DefaultProfile: "work",
		Profiles: map[string]Profile{
			"work": {URL: "https://redmine.company.com", APIKey: "abc123"},
		},
	}

	_, err := cfg.GetProfile("unknown")
	if err == nil {
		t.Fatal("GetProfile with unknown name should return error")
	}
}

func TestLoadConfig_URLTrailingSlashRemoval(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `default_profile: work
profiles:
  work:
    url: https://redmine.company.com/
    api_key: abc123
  oss:
    url: https://redmine.example.org///
    api_key: def456
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	work := cfg.Profiles["work"]
	if work.URL != "https://redmine.company.com" {
		t.Fatalf("trailing slash not removed for work: got %q", work.URL)
	}

	oss := cfg.Profiles["oss"]
	if oss.URL != "https://redmine.example.org" {
		t.Fatalf("trailing slashes not removed for oss: got %q", oss.URL)
	}
}
