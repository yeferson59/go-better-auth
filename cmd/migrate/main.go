package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	PostgresDriver = "postgres"
	SQLiteDriver   = "sqlite3"
)

func main() {
	var (
		driver    = flag.String("driver", PostgresDriver, "Database driver (postgres, sqlite3)")
		dsn       = flag.String("dsn", "", "Database DSN")
		direction = flag.String("direction", "up", "Migration direction (up, down)")
		steps     = flag.Int("steps", 0, "Number of migration steps (0 for all)")
		version   = flag.Int("version", 0, "Migrate to specific version")
		force     = flag.Int("force", 0, "Force migration to specific version")
		create    = flag.String("create", "", "Create new migration file")
	)
	flag.Parse()

	if *dsn == "" && *create == "" {
		log.Fatal("DSN is required")
	}

	// Handle create migration
	if *create != "" {
		if err := createMigration(*create, *driver); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Setup database connection
	db, err := sql.Open(*driver, *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get migration path
	migrationPath := getMigrationPath(*driver)

	// Setup migrate instance
	m, err := setupMigrate(db, *driver, migrationPath)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	// Execute migration commands
	switch {
	case *force != 0:
		if err := m.Force(*force); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Forced migration to version %d\n", *force)

	case *version != 0:
		if err := m.Migrate(uint(*version)); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Migrated to version %d\n", *version)

	case *direction == "up":
		if *steps == 0 {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatal(err)
			}
			fmt.Println("Migration up completed")
		} else {
			if err := m.Steps(*steps); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Migration up %d steps completed\n", *steps)
		}

	case *direction == "down":
		if *steps == 0 {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatal(err)
			}
			fmt.Println("Migration down completed")
		} else {
			if err := m.Steps(-*steps); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Migration down %d steps completed\n", *steps)
		}

	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}

func setupMigrate(db *sql.DB, driver, migrationPath string) (*migrate.Migrate, error) {
	var dbDriver database.Driver
	var err error

	switch driver {
	case PostgresDriver:
		dbDriver, err = postgres.WithInstance(db, &postgres.Config{})
	case SQLiteDriver:
		dbDriver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err != nil {
		return nil, err
	}

	sourceDriver, err := (&file.File{}).Open(migrationPath)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("file", sourceDriver, driver, dbDriver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func getMigrationPath(driver string) string {
	switch driver {
	case PostgresDriver:
		return "file://internal/infra/db/migrations/postgres"
	case SQLiteDriver:
		return "file://internal/infra/db/migrations/sqlite"
	default:
		return "file://internal/infra/db/migrations/postgres"
	}
}

func createMigration(name, driver string) error {
	migrationDir := fmt.Sprintf("internal/infra/db/migrations/%s", driver)
	if driver == PostgresDriver {
		migrationDir = "internal/infra/db/migrations/postgres"
	}

	// Get next migration number
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	nextNum := 1
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Parse migration number from filename
		var num int
		if _, err := fmt.Sscanf(file.Name(), "%d_", &num); err == nil {
			if num >= nextNum {
				nextNum = num + 1
			}
		}
	}

	// Create migration files
	upFile := fmt.Sprintf("%s/%03d_%s.up.sql", migrationDir, nextNum, name)
	downFile := fmt.Sprintf("%s/%03d_%s.down.sql", migrationDir, nextNum, name)

	// Create up migration
	if err := os.WriteFile(upFile, fmt.Appendf([]byte{}, "-- Migration: %s\n-- Up migration\n", name), 0644); err != nil {
		return err
	}

	// Create down migration
	if err := os.WriteFile(downFile, fmt.Appendf([]byte{}, "-- Migration: %s\n-- Down migration\n", name), 0644); err != nil {
		return err
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upFile)
	fmt.Printf("  %s\n", downFile)

	return nil
}
