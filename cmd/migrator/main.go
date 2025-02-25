package main

import (
	"errors"
	"flag"
	"fmt"

	"restrocheck/internal/config"
	"restrocheck/internal/storage/postgresql"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	var migrationsPath string
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	cfg := config.MustLoad()
	var dbURL string = postgresql.BuildStringConnectDB(cfg)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied successfull")
}
