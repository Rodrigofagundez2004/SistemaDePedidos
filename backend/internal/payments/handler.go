package payments

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RodrigoFagundez2004/SistemaDePedidos/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(db *gorm.DB, ordersClient OrdersClient) *Handler {
	repo := NewRepository(db)
	service := NewService(repo, ordersClient)
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, validator middleware.TokenValidator) {
	group := rg.Group("/payments")
	group.Use(middleware.RequireAuth(validator))
	{
		group.POST("", h.Create)
		group.GET("/:id", h.GetByID)
		group.GET("/order/:orderID", h.ListByOrder)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreatePayment(&req, c.GetHeader("Authorization"))
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrOrderAlreadyPaid),
			errors.Is(err, ErrOrderNotPending),
			errors.Is(err, ErrIdempotencyKeyConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrOrdersServiceUnavailable):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		default:
			slog.Error("error inesperado en payments", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	payment, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToPaymentResponse(payment))
}

func (h *Handler) ListByOrder(c *gin.Context) {
	orderID := c.Param("orderID")

	payments, err := h.service.ListByOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToPaymentResponseList(payments))
}
