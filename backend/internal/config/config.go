package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	AppEnv            string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBSSLMode         string
	UsersDBName       string
	ProductsDBName    string
	OrdersDBName      string
	PaymentsDBName    string
	JWTSecret         string
	JWTExpirationHour int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		expHours = 24
	}

	cfg := &Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		AppEnv:            getEnv("APP_ENV", "development"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		UsersDBName:       getEnv("USERS_DB_NAME", "users_db"),
		ProductsDBName:    getEnv("PRODUCTS_DB_NAME", "products_db"),
		OrdersDBName:      getEnv("ORDERS_DB_NAME", "orders_db"),
		PaymentsDBName:    getEnv("PAYMENTS_DB_NAME", "payments_db"),
		JWTSecret:         getEnv("JWT_SECRET", "dev_secret"),
		JWTExpirationHour: expHours,
	}

	return cfg, nil
}

func (c *Config) DSN(dbName string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, dbName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
