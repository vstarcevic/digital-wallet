package messaging

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"strconv"
	"time"
	"transaction-service/api"
	"transaction-service/database"
	"transaction-service/model"

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

func SubscribeNats(cfg api.Config) {
	cfg.Nts.Subscribe("balance", func(m *nats.Msg) {
		getBalance(cfg.Nts, m, cfg.Db)
	})
}

func getBalance(nc *nats.Conn, m *nats.Msg, db *sql.DB) (*int, error) {

	userId, err := strconv.Atoi(string(m.Data))
	if err != nil {
		return nil, err
	}

	balance, err := database.GetBalance(db, userId)

	userResponse := model.UserBalanceResponse{
		UserId:  &userId,
		Balance: *balance,
		Error:   nil,
	}

	out, err := json.Marshal(userResponse)
	if err != nil {
		return nil, err
	}

	nc.Publish(m.Reply, out)

	return &userId, nil
}
