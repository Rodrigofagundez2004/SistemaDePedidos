package database

import (
	"log/slog"

	"github.com/RodrigoFagundez2004/SistemaDePedidos/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connections struct {
	Users    *gorm.DB
	Products *gorm.DB
	Orders   *gorm.DB
	Payments *gorm.DB
}

func NewConnections(cfg *config.Config) (*Connections, error) {
	users, err := connect(cfg.DSN(cfg.UsersDBName), "users_db")
	if err != nil {
		return nil, err
	}
	products, err := connect(cfg.DSN(cfg.ProductsDBName), "products_db")
	if err != nil {
		return nil, err
	}
	orders, err := connect(cfg.DSN(cfg.OrdersDBName), "orders_db")
	if err != nil {
		return nil, err
	}
	payments, err := connect(cfg.DSN(cfg.PaymentsDBName), "payments_db")
	if err != nil {
		return nil, err
	}

	slog.Info("conexiones a bases de datos establecidas")
	return &Connections{
		Users:    users,
		Products: products,
		Orders:   orders,
		Payments: payments,
	}, nil
}

func connect(dsn, name string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		slog.Error("error conectando a base de datos", "db", name, "error", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		slog.Error("ping a base de datos falló", "db", name, "error", err)
		return nil, err
	}

	slog.Info("conectado", "db", name)
	return db, nil
}
