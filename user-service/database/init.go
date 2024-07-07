package database

import (
	"database/sql"
	"log"

	"log/slog"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectToDB(dsn string) *sql.DB {
	counts := 0

	for {
		connection, err := openDB(dsn)

		if err != nil {
			slog.Warn("Postgres not yet ready ...")
			counts++
		} else {
			slog.Info("Connected to postgres!")
			return connection
		}

		if counts > 10 {
			slog.Error("Cannot connect to postgres.")
			log.Panic(err)
			return nil
		}

		slog.Warn("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
