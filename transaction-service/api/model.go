package api

import (
	"database/sql"

	"github.com/nats-io/nats.go"
)

type Config struct {
	Db  *sql.DB
	Nts *nats.Conn
	App AppSettings
}

type AppSettings struct {
	Dsn      string
	NatsUrl  string
	KafkaUrl string
}
