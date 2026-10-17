package tests

import (
	"strings"
	"testing"

	"github.com/RodrigoFagundez2004/SistemaDePedidos/backend/internal/config"
)

func TestDSNTieneTodosLosDatos(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "secreto",
		DBSSLMode:  "disable",
	}

	dsn := cfg.DSN("users_db")

	esperados := []string{
		"host=localhost",
		"port=5432",
		"user=postgres",
		"password=secreto",
		"dbname=users_db",
		"sslmode=disable",
	}
	for _, parte := range esperados {
		if !strings.Contains(dsn, parte) {
			t.Errorf("el DSN %q no contiene %q", dsn, parte)
		}
	}
}

func TestLoadUsaValoresPorDefecto(t *testing.T) {
	// Una variable vacía se trata como no configurada, así que se usa el valor por defecto.
	t.Setenv("APP_PORT", "")
	t.Setenv("JWT_EXPIRATION_HOURS", "")
	t.Setenv("DB_HOST", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.AppPort != "8080" {
		t.Errorf("AppPort = %s, want 8080", cfg.AppPort)
	}
	if cfg.JWTExpirationHour != 24 {
		t.Errorf("JWTExpirationHour = %d, want 24", cfg.JWTExpirationHour)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost = %s, want localhost", cfg.DBHost)
	}
}
