package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Изолируем от .env файлов в проекте — chdir во временную директорию.
	tmp := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	cfg := &Config{}
	if err := Load(cfg, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.PgHost != "localhost" {
		t.Errorf("expected PgHost 'localhost', got %q", cfg.PgHost)
	}
	if cfg.PgPort != "5432" {
		t.Errorf("expected PgPort '5432', got %q", cfg.PgPort)
	}
	if cfg.PgUser != "postgres" {
		t.Errorf("expected PgUser 'postgres', got %q", cfg.PgUser)
	}
	if cfg.PgPassword != "postgres" {
		t.Errorf("expected PgPassword 'postgres', got %q", cfg.PgPassword)
	}
	if cfg.PgDBName != "store_scraper" {
		t.Errorf("expected PgDBName 'store_scraper', got %q", cfg.PgDBName)
	}
	if cfg.RedisHost != "localhost" {
		t.Errorf("expected RedisHost 'localhost', got %q", cfg.RedisHost)
	}
	if cfg.RedisPort != "6379" {
		t.Errorf("expected RedisPort '6379', got %q", cfg.RedisPort)
	}
	if cfg.HTTPPort != "8080" {
		t.Errorf("expected HTTPPort '8080', got %q", cfg.HTTPPort)
	}
	if cfg.GinMode != "debug" {
		t.Errorf("expected GinMode 'debug', got %q", cfg.GinMode)
	}
	if cfg.RateLimitRPS != 100 {
		t.Errorf("expected RateLimitRPS 100, got %d", cfg.RateLimitRPS)
	}
	if cfg.ScrapeInterval != 10*time.Minute {
		t.Errorf("expected ScrapeInterval 10m, got %v", cfg.ScrapeInterval)
	}
	if cfg.ScrapeWorkers != 5 {
		t.Errorf("expected ScrapeWorkers 5, got %d", cfg.ScrapeWorkers)
	}
	if cfg.LogMode != "dev" {
		t.Errorf("expected LogMode 'dev', got %q", cfg.LogMode)
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	cfg := &Config{}
	err := Load(cfg, "/nonexistent/path/config.env")
	if err == nil {
		t.Fatal("expected error for invalid file, got nil")
	}
}

func TestRedisConfig_Addr(t *testing.T) {
	cfg := RedisConfig{
		RedisHost: "redis-server",
		RedisPort: "6380",
	}

	addr := cfg.Addr()
	if addr != "redis-server:6380" {
		t.Errorf("expected 'redis-server:6380', got %q", addr)
	}
}

func TestPostgresConfig_Getters(t *testing.T) {
	cfg := PostgresConfig{
		PgHost:     "db-host",
		PgPort:     "5433",
		PgUser:     "admin",
		PgPassword: "secret",
		PgDBName:   "mydb",
	}

	if cfg.GetHost() != "db-host" {
		t.Errorf("GetHost: expected 'db-host', got %q", cfg.GetHost())
	}
	if cfg.GetPort() != "5433" {
		t.Errorf("GetPort: expected '5433', got %q", cfg.GetPort())
	}
	if cfg.GetUser() != "admin" {
		t.Errorf("GetUser: expected 'admin', got %q", cfg.GetUser())
	}
	if cfg.GetPassword() != "secret" {
		t.Errorf("GetPassword: expected 'secret', got %q", cfg.GetPassword())
	}
	if cfg.GetDBName() != "mydb" {
		t.Errorf("GetDBName: expected 'mydb', got %q", cfg.GetDBName())
	}
}

func TestBaseConfig_EnvChecks(t *testing.T) {
	dev := BaseConfig{Environment: EnvDev}
	if !dev.IsDevEnv() {
		t.Error("expected IsDevEnv true for dev")
	}
	if dev.IsProdEnv() {
		t.Error("expected IsProdEnv false for dev")
	}

	prod := BaseConfig{Environment: EnvProd}
	if !prod.IsProdEnv() {
		t.Error("expected IsProdEnv true for prod")
	}
	if prod.IsDevEnv() {
		t.Error("expected IsDevEnv false for prod")
	}

	empty := BaseConfig{}
	if empty.IsDevEnv() || empty.IsProdEnv() {
		t.Error("expected both false for empty environment")
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	// Изолируем от .env файлов.
	tmp := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	t.Setenv("PG_HOST", "custom-host")
	t.Setenv("PG_PORT", "9999")
	t.Setenv("HTTP_PORT", "3000")

	cfg := &Config{}
	if err := Load(cfg, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.PgHost != "custom-host" {
		t.Errorf("expected PgHost 'custom-host', got %q", cfg.PgHost)
	}
	if cfg.PgPort != "9999" {
		t.Errorf("expected PgPort '9999', got %q", cfg.PgPort)
	}
	if cfg.HTTPPort != "3000" {
		t.Errorf("expected HTTPPort '3000', got %q", cfg.HTTPPort)
	}
}
