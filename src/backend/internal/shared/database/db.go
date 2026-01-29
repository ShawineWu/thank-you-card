package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/castlery/thank-you-card/internal/card/domain/models"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Initialize sets up the database connection
func Initialize(config Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=public",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("Database connection established successfully")

	// Diagnostics: Check current database and tables
	var dbName string
	db.Raw("SELECT current_database()").Scan(&dbName)
	log.Printf("Connected to database: %s", dbName)

	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	log.Printf("Tables in public schema: %v", tables)

	return nil
}

// AutoMigrate handles automatic schema migration
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Running database auto-migration...")
	err := DB.AutoMigrate(
		&models.Card{},
		&models.CardRecipient{},
		&models.CardValue{},
		&models.CompanyValue{},
		&models.EmployeeMilestone{},
	)
	if err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Println("Database auto-migration completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

// HealthCheck performs a database health check
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
