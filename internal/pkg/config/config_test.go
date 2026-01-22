package config

import (
	"os"
	"testing"
)

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("DATABASE_DSN", "root:root@tcp(localhost:3306)/myapp?parseTime=true")
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("TEMPLATES_PATH", "templates/*.html")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.Database.DSN == "" {
		t.Fatalf("expected database DSN to be set")
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("expected host 127.0.0.1, got %q", cfg.Server.Host)
	}

	if cfg.Server.Port != 9090 {
		t.Fatalf("expected port 9090, got %d", cfg.Server.Port)
	}

	if cfg.Templates.Path != "templates/*.html" {
		t.Fatalf("expected templates path, got %q", cfg.Templates.Path)
	}
}

func TestLoad_EmptyDSN_ReturnsError(t *testing.T) {
	os.Unsetenv("DATABASE_DSN")

	cfg, err := Load()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if cfg != nil {
		t.Fatalf("expected nil config on error")
	}
}
