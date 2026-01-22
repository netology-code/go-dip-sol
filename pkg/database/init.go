package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Config содержит конфигурацию базы данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(cfg Config) (*sql.DB, error) {

	// 1. Сформировать строку подключения (DSN)
	dsn := GetDSN(cfg)

	// 2. Открыть соединение с БД
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 3. Проверить соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 4. Настроить пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to PostgreSQL database")
	return db, nil
}

// GetDSN формирует строку подключения к PostgreSQL
func GetDSN(cfg Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

// Close закрывает соединение с базой данных
func Close(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}

// TestConnection выполняет тестовый запрос к БД
func TestConnection(db *sql.DB) error {
	var result int
	return db.QueryRow("SELECT 1").Scan(&result)
}
