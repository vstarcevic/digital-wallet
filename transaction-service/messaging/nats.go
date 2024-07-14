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

	var errorText string

	userId, err := strconv.Atoi(string(m.Data))
	if err != nil {
		errorText = "internal error"
	}

	balance, err := database.GetBalance(db, userId)
	if err != nil {
		errorText = "user not found"
	}

	if errorText != "" {
		out, _ := json.Marshal(model.UserBalanceResponse{Error: &errorText})
		nc.Publish(m.Reply, out)
		return nil, err
	}

	userResponse := model.UserBalanceResponse{
		UserId:  &userId,
		Balance: *balance,
		Error:   nil,
	}

	out, _ := json.Marshal(userResponse)

	nc.Publish(m.Reply, out)

	return &userId, nil
}
