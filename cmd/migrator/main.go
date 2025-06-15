package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string
	flag.StringVar(&storagePath, "storage-path", "", "path to storage file")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to the migrations directory")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of the migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migrationsPath is required")
	}

	if migrationsTable == "" {
		panic("migrationsTable is required")
	}

	m, err := migrate.New("file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable))
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no change to apply")
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
