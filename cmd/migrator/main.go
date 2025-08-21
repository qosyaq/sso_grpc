package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		storagePath     string
		postgresDSN     string
		migrationsPath  string
		migrationsTable string
		down            bool
		steps           int
	)

	flag.StringVar(&storagePath, "storage-path", "", "path to SQLite storage (if used)")
	flag.StringVar(&postgresDSN, "postgres-dsn", "", "Postgres DSN (if used)")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.BoolVar(&down, "down", false, "if set, run migrations down")
	flag.IntVar(&steps, "steps", 1, "number of steps to migrate down (ignored if --down not set)")
	flag.Parse()

	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	var dbURL string

	if postgresDSN != "" {
		// Postgres
		dbURL = fmt.Sprintf("%s&x-migrations-table=%s", postgresDSN, migrationsTable)
	} else if storagePath != "" {
		// SQLite
		dbURL = fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable)
	} else {
		log.Fatal("either --postgres-dsn or --storage-path must be provided")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	if down {
		if err := m.Steps(-steps); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to rollback")
				return
			}
			log.Fatal(err)
		}
		fmt.Printf("rolled back %d migration(s)\n", steps)
		return
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatal(err)
	}

	fmt.Println("migrations applied")
}
