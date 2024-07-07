package database

import (
	"log"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(dsn string) {
	m, err := migrate.New(
		"file://../../database/migrations",
		dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		slog.Warn("Migration", "Error:", err)
	}
}
