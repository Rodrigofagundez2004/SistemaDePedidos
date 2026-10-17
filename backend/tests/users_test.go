package tests

import (
	"testing"
	"time"

	"github.com/RodrigoFagundez2004/SistemaDePedidos/backend/internal/users"
	jwt "github.com/golang-jwt/jwt/v5"
)

// firmar crea un JWT firmado con el secreto indicado, como lo hace el servicio de usuarios.
func firmar(t *testing.T, secreto string, claims jwt.MapClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secreto))
	if err != nil {
		t.Fatalf("firmar token: %v", err)
	}
	return token
}

func TestValidateTokenAceptaTokenValido(t *testing.T) {
	svc := users.NewService(nil, "secreto-de-prueba", 1)
	token := firmar(t, "secreto-de-prueba", jwt.MapClaims{
		"sub":   "user-123",
		"email": "ana@x.com",
		"role":  "customer",
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("el token válido fue rechazado: %v", err)
	}
	if (*claims)["sub"] != "user-123" {
		t.Errorf("sub = %v, want user-123", (*claims)["sub"])
	}
	if (*claims)["role"] != "customer" {
		t.Errorf("role = %v, want customer", (*claims)["role"])
	}
}

func TestValidateTokenRechazaOtroSecreto(t *testing.T) {
	svc := users.NewService(nil, "secreto-del-servidor", 1)
	token := firmar(t, "secreto-de-un-atacante", jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	if _, err := svc.ValidateToken(token); err == nil {
		t.Error("un token firmado con otro secreto no debería ser válido")
	}
}

func TestValidateTokenRechazaTokenExpirado(t *testing.T) {
	svc := users.NewService(nil, "secreto-de-prueba", 1)
	token := firmar(t, "secreto-de-prueba", jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})

	if _, err := svc.ValidateToken(token); err == nil {
		t.Error("un token expirado no debería ser válido")
	}
}
