package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenValidator interface {
	ValidateToken(token string) (*jwt.MapClaims, error)
}

func RequireAuth(validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "falta el header Authorization",
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "formato inválido. Usá: Bearer <token>",
			})
			return
		}

		claims, err := validator.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido o expirado",
			})
			return
		}

		// Guardar datos del usuario en el contexto
		c.Set("user_id", (*claims)["sub"])
		c.Set("user_email", (*claims)["email"])
		c.Set("user_role", (*claims)["role"])

		c.Next()
	}
}
