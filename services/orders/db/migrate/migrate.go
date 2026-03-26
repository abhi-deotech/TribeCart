package migrate

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

//go:embed ../../migrations/*.sql
var migrationsFS embed.FS

// Migration represents a database migration with version, name, and SQL content.
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// MigrationError represents an error that occurred during migration.
type MigrationError struct {
	Migration *Migration
	Err       error
}

func (e *MigrationError) Error() string {
	return fmt.Sprintf("migration %s failed: %v", e.Migration.Name, e.Err)
}

// RunMigrations runs all pending migrations up to the latest version.
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migrations from the embedded filesystem
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get current database version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current database version: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Apply pending migrations
	for _, m := range migrations {
		if m.Version > currentVersion {
			log.Printf("Applying migration %s (version %d)", m.Name, m.Version)
			
			// Start transaction
			tx, err := db.Begin()
			if err != nil {
				return &MigrationError{m, fmt.Errorf("failed to begin transaction: %w", err)}
			}

			// Execute migration
			if _, err := tx.Exec(m.Up); err != nil {
				tx.Rollback()
				return &MigrationError{m, fmt.Errorf("failed to apply migration: %w", err)}
			}

			// Update schema version
			if _, err := tx.Exec(
				"INSERT INTO schema_migrations (version, name) VALUES ($1, $2) ON CONFLICT (version) DO UPDATE SET name = EXCLUDED.name, applied_at = NOW()",
				m.Version,
				m.Name,
			); err != nil {
				tx.Rollback()
				return &MigrationError{m, fmt.Errorf("failed to update schema version: %w", err)}
			}

			// Commit transaction
			if err := tx.Commit(); err != nil {
				return &MigrationError{m, fmt.Errorf("failed to commit transaction: %w", err)}
			}

			log.Printf("Successfully applied migration %s (version %d)", m.Name, m.Version)
		}
	}

	return nil
}

// RollbackMigrations rolls back the most recent migration.
func RollbackMigrations(db *sql.DB) error {
	// Get current database version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current database version: %w", err)
	}

	if currentVersion == 0 {
		return errors.New("no migrations to rollback")
	}

	// Load all migrations
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Find the migration to rollback
	var target *Migration
	for _, m := range migrations {
		if m.Version == currentVersion {
			target = m
			break
		}
	}

	if target == nil {
		return fmt.Errorf("no migration found for version %d", currentVersion)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute rollback
	if _, err := tx.Exec(target.Down); err != nil {
		tx.Rollback()
		return &MigrationError{target, fmt.Errorf("failed to rollback migration: %w", err)}
	}

	// Update schema version to the previous one
	var prevVersion int
	err = tx.QueryRow(
		"SELECT COALESCE(MAX(version), 0) FROM schema_migrations WHERE version < $1",
		currentVersion,
	).Scan(&prevVersion)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return fmt.Errorf("failed to get previous migration version: %w", err)
	}

	// Delete the current version
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", currentVersion); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update schema version: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully rolled back migration %s (version %d)", target.Name, target.Version)
	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist.
func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

// getCurrentVersion returns the current database schema version.
func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return version, nil
}

// loadMigrations loads all migration files from the embedded filesystem.
func loadMigrations() ([]*Migration, error) {
	var migrations []*Migration
	migrationFiles := make(map[int]struct{})

	// Read all files from the migrations directory
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Process each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Parse filename (e.g., 001_initial_schema.up.sql)
		name := file.Name()
		parts := strings.Split(name, "_")
		if len(parts) < 3 {
			continue
		}

		// Extract version number
		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid migration filename format: %s", name)
		}

		// Skip if we've already processed this version
		if _, exists := migrationFiles[version]; exists {
			continue
		}
		migrationFiles[version] = struct{}{}

		// Extract migration name (e.g., "initial_schema")
		migrationName := strings.TrimSuffix(
			strings.TrimPrefix(name, fmt.Sprintf("%03d_", version)),
			filepath.Ext(name),
		)
		migrationName = strings.TrimSuffix(migrationName, ".up")

		// Read up migration
		upSQL, err := fs.ReadFile(migrationsFS, filepath.Join("migrations", name))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		// Read down migration
		downFile := strings.Replace(name, ".up.sql", ".down.sql", 1)
		downSQL, err := fs.ReadFile(migrationsFS, filepath.Join("migrations", downFile))
		if err != nil {
			return nil, fmt.Errorf("failed to read down migration file %s: %w", downFile, err)
		}

		// Create and add migration
		migrations = append(migrations, &Migration{
			Version: version,
			Name:    migrationName,
			Up:      string(upSQL),
			Down:    string(downSQL),
		})
	}

	return migrations, nil
}

// ValidateMigrations checks if all migrations are valid.
func ValidateMigrations() error {
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Check for duplicate versions
	versions := make(map[int]bool)
	for _, m := range migrations {
		if versions[m.Version] {
			return fmt.Errorf("duplicate migration version: %d", m.Version)
		}
		versions[m.Version] = true

		// Validate migration name format
		if !regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(m.Name) {
			return fmt.Errorf("invalid migration name: %s (must contain only lowercase letters, numbers, and underscores)", m.Name)
		}

		// Validate up migration is not empty
		if strings.TrimSpace(m.Up) == "" {
			return fmt.Errorf("empty up migration for version %d", m.Version)
		}

		// Validate down migration is not empty
		if strings.TrimSpace(m.Down) == "" {
			return fmt.Errorf("empty down migration for version %d", m.Version)
		}
	}

	// Ensure versions are sequential
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for i, m := range migrations {
		if m.Version != i+1 {
			return fmt.Errorf("missing migration: expected version %d, got %d", i+1, m.Version)
		}
	}

	return nil
}
