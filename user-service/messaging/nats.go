package messaging

import (
	"log"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

func ConnectToNats(url string) *nats.Conn {

	counts := 0

	for {
		connection, err := nats.Connect(url)

		if err != nil {
			slog.Warn("Nats not yet ready ...")
			counts++
		} else {
			slog.Info("Connected to Nats!")
			return connection
		}

		if counts > 10 {
			slog.Error("Cannot connect to nats.")
			log.Panic(err)
			return nil
		}

		slog.Warn("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

}
