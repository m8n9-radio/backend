package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Migrate struct {
	m *migrate.Migrate
}

func NewMigrate(dns string) *Migrate {
	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Migrate{
		m: m,
	}
}

func CreateMigrationFile(filepath, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func MigrateCreate(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join("migrations", filename+".up.sql")
	downFile := filepath.Join("migrations", filename+".down.sql")

	if err := CreateMigrationFile(upFile, "-- Migration up"); err != nil {
		return fmt.Errorf("Error creating up migration file: %v\n", err)
	}

	if err := CreateMigrationFile(downFile, "-- Migration down"); err != nil {
		return fmt.Errorf("Error creating down migration file: %v\n", err)
	}

	fmt.Printf("Created migrations:\n\t%s\n\t%s\n", upFile, downFile)

	return nil
}

func (m *Migrate) Up(steps int) error {
	if steps > 0 {
		if err := m.m.Steps(steps); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No migrations to run")
				return nil
			}
			return fmt.Errorf("failed to run %d migration steps: %w", steps, err)
		}
		fmt.Printf("Successfully ran %d migration steps\n", steps)
	} else {
		if err := m.m.Up(); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No migrations to run")
				return nil
			}
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		fmt.Println("Migrations completed successfully")
	}

	return nil
}

func (m *Migrate) Down(steps int) error {
	if steps > 0 {
		if err := m.m.Steps(-steps); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No migrations to rollback")
				return nil
			}
			return fmt.Errorf("failed to rollback %d migration steps: %w", steps, err)
		}
		fmt.Printf("Successfully rolled back %d migration steps\n", steps)
	} else {
		if err := m.m.Down(); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No migrations to rollback")
				return nil
			}
			return fmt.Errorf("failed to rollback all migrations: %w", err)
		}
		fmt.Println("All migrations rolled back successfully")
	}

	return nil
}

func (m *Migrate) Drop() error {
	if err := m.m.Drop(); err != nil {
		return fmt.Errorf("failed to drop database schema: %w", err)
	}

	fmt.Println("Database schema dropped successfully")
	return nil
}

func (m *Migrate) Force(version int) error {
	if err := m.m.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version to %d: %w", version, err)
	}

	fmt.Printf("Forced migration version to %d\n", version)
	return nil
}

func (m *Migrate) Status() error {
	version, dirty, err := m.m.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	fmt.Printf("Current migration version: %d\n", version)
	if dirty {
		fmt.Println("Status: DIRTY (migration failed)")
	} else {
		fmt.Println("Status: CLEAN")
	}

	return nil
}
