package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RodrigoFagundez2004/SistemaDePedidos/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// fakeValidator simula el servicio de usuarios: solo acepta el token "valido".
type fakeValidator struct{}

func (fakeValidator) ValidateToken(token string) (*jwt.MapClaims, error) {
	if token != "valido" {
		return nil, http.ErrNoCookie // cualquier error sirve para el test
	}
	return &jwt.MapClaims{"sub": "user-123", "email": "ana@x.com", "role": "customer"}, nil
}

// newRouter arma una ruta protegida de prueba.
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/privado", middleware.RequireAuth(fakeValidator{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": c.GetString("user_id")})
	})
	return r
}

func TestRequireAuth(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{"sin header Authorization", "", http.StatusUnauthorized},
		{"formato incorrecto", "Basic valido", http.StatusUnauthorized},
		{"token inválido", "Bearer invalido", http.StatusUnauthorized},
		{"token válido", "Bearer valido", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/privado", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			newRouter().ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestRequireAuthGuardaElUsuarioEnElContexto(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/privado", nil)
	req.Header.Set("Authorization", "Bearer valido")
	w := httptest.NewRecorder()

	newRouter().ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "user-123") {
		t.Errorf("la respuesta no tiene el user_id; body = %s", w.Body.String())
	}
}
