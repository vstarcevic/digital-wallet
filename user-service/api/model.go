package api

import "database/sql"

type Config struct {
	Db *sql.DB
}
