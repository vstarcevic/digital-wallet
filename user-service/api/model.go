package api

import (
	"database/sql"

	"github.com/nats-io/nats.go"
)

type Config struct {
	Db  *sql.DB
	Tx  *sql.Tx
	Nts *nats.Conn
}
