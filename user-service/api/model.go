package api

import (
	"database/sql"

	"github.com/nats-io/nats.go"
)

type Config struct {
	Db  *sql.DB
	Nts *nats.Conn
}
