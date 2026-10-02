package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/config"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/database"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/middleware"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/orders"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/payments"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/products"
	"github.com/rf4tc/sistema-de-pedidos/backend/internal/users"
)

func main() {
	// ---------------------------------------------------------
	// 1. Configurar logs estructurados
	// ---------------------------------------------------------
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// ---------------------------------------------------------
	// 2. Cargar configuración desde .env
	// ---------------------------------------------------------
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error cargando configuración", "error", err)
		os.Exit(1)
	}
	slog.Info("configuración cargada", "env", cfg.AppEnv, "port", cfg.AppPort)

	// ---------------------------------------------------------
	// 3. Conectar a las 4 bases de datos
	// ---------------------------------------------------------
	conns, err := database.NewConnections(cfg)
	if err != nil {
		slog.Error("error conectando a las bases de datos", "error", err)
		os.Exit(1)
	}

	// ---------------------------------------------------------
	// 4. Crear el router de Gin
	// ---------------------------------------------------------
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// ---------------------------------------------------------
	// 5. Health check general
	// ---------------------------------------------------------
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"env":    cfg.AppEnv,
		})
	})

	// ---------------------------------------------------------
	// 6. Registrar handlers por dominio
	// ---------------------------------------------------------
	api := router.Group("/api/v1")

	users.NewHandler(conns.Users).RegisterRoutes(api)
	products.NewHandler(conns.Products).RegisterRoutes(api)
	orders.NewHandler(conns.Orders).RegisterRoutes(api)
	payments.NewHandler(conns.Payments).RegisterRoutes(api)

	// ---------------------------------------------------------
	// 7. Servidor HTTP con graceful shutdown
	// ---------------------------------------------------------
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		slog.Info("servidor arrancando", "port", cfg.AppPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error en el servidor", "error", err)
			os.Exit(1)
		}
	}()

	// ---------------------------------------------------------
	// 8. Esperar señal de cierre (Ctrl+C)
	// ---------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("apagando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error al apagar", "error", err)
	}
	slog.Info("servidor apagado correctamente")
}
