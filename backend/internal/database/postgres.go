package database

import (
	"log/slog"
	"https://github.com/Rodrigofagundez2004/SistemaDePedidos/blob/main/backend/internal/config/config.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connections struct {
	UsersDB    *gorm.DB
	ProductsDB *gorm.DB
	OrdersDB   *gorm.DB
	PaymentsDB *gorm.DB
}
func NewConnections(cfg *config.Config) (*Connections, error) {
	users, err := connect(cfg.DSN(cfg.UsersDBName),"users_db")
	if err != nil {
		return nil, err;
	}
	products, err := connect(cfg.DSN(cfg.ProductsDBName),"products_db")
	if err != nil {
		return nil, err;
	}
	orders, err := connect(cfg.DSN(cfg.OrdersDBName),"orders_db")
	if err != nil {
		return nil, err;
	}
	payments, err := connect(cfg.DSN(cfg.PaymentsDBName),"payments_db")
	if err != nil 
	{
		return nil, err;
	}
	slog.Info("Connected to all databases successfully")
	return &Connections{
		UsersDB:    users,
		ProductsDB: products,
		OrdersDB:   orders,
		PaymentsDB: payments,
	}, nil
}
func connect(dsn, name string) (*gorm.DB, error)
{
	db, err : = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil
	{
		slog.Error("Error for conection to database "+name, "error", err)
		return nil, err

	}
	sqlDB, err := db.DB()
	if err != nil 
	{
		return nil, err

	}
	if err := sqlDB.Ping(); err != nil
	{
		slog.Error("Error pinging database "+name, "error", err)
		return nil, err
	}
	slog.Info("conectado", "db", name)
	return db, nil
}