package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {

	return &Handler{db: db}
}
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/users")
	{
		group.GET("/ping", h.Ping)
	}
}
func (h *Handler) Ping(c *gin.Context) {

	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ping database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"service": "users",
		"status":  "ok",
		"db":      "connected",
	})
}
