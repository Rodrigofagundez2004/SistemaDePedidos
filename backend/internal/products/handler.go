package products

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
	group := rg.Group("/products")
	{
		group.GET("/ping", h.Ping)
	}
}

// la idea del ping es que verficica si la conexion a mi db funciona de prodcuts en este caso
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
		"service": "products",
		"status":  "ok",
		"db":      "connected",
	})
}
