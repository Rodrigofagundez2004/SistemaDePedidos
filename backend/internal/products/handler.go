package products

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		group.GET("", h.List)
		group.GET("/ping", h.Ping)
		group.GET("/:id", h.GetByID)
	}
}

// List devuelve los productos activos: GET /api/v1/products?page=1&limit=20&search=...
func (h *Handler) List(c *gin.Context) {
	var q ListProductsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parámetros inválidos"})
		return
	}

	service := NewService(&Repository{db: h.db})
	items, total, err := service.ListProducts(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al listar productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

// GetByID devuelve un producto: GET /api/v1/products/:id (lo usa el servicio de Orders)
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	service := NewService(&Repository{db: h.db})
	product, err := service.GetProductById(id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar producto"})
		return
	}

	c.JSON(http.StatusOK, product)
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
