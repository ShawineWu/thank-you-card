package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// MigrationRecord tracks which migrations have been applied
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Filename  string `gorm:"unique;not null"`
	AppliedAt string `gorm:"not null"`
}

// RunSQLMigrations runs SQL migration files from the migrations directory
func RunSQLMigrations(db *gorm.DB) error {
	// Create migrations table if it doesn't exist
	err := db.AutoMigrate(&MigrationRecord{})
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get already applied migrations
	var appliedMigrations []MigrationRecord
	db.Find(&appliedMigrations)
	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Filename] = true
	}

	// Apply pending migrations
	for _, filename := range migrationFiles {
		if appliedMap[filename] {
			log.Printf("Migration %s already applied, skipping...", filename)
			continue
		}

		log.Printf("Applying migration: %s", filename)
		err := applyMigration(db, filename)
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", filename, err)
		}

		// Record the migration as applied
		migrationRecord := MigrationRecord{
			Filename:  filename,
			AppliedAt: "NOW()",
		}
		db.Create(&migrationRecord)
		log.Printf("Successfully applied migration: %s", filename)
	}

	return nil
}

// getMigrationFiles returns a sorted list of migration files
func getMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Sort files to ensure they run in order
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// applyMigration applies a single migration file
func applyMigration(db *gorm.DB, filename string) error {
	migrationPath := filepath.Join("migrations", filename)
	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	// Split the content by semicolons to handle multiple statements
	statements := strings.Split(string(content), ";")
	
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		err := db.Exec(statement).Error
		if err != nil {
			return fmt.Errorf("failed to execute statement in %s: %w\nStatement: %s", filename, err, statement)
		}
	}

	return nil
}