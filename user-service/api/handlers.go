package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"user-service/database"
	"user-service/messaging"

	m "user-service/model"
)

func (cfg *Config) createUser(w http.ResponseWriter, r *http.Request) {

	var requestPayload m.JsonRequest
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		writeError(w, http.StatusNotAcceptable, errors.New("error unmarshaling request"))
		return
	}

	if requestPayload.Email == "" {
		writeError(w, http.StatusNotAcceptable, errors.New("email empty"))
		return
	}

	ctx := context.Background()

	tx, err := cfg.Db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp, err := database.CreateUserWithTx(ctx, tx, cfg.Db, requestPayload.Email)
	if err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// create in kafka
	err = messaging.PublishJSON("user-created", resp, cfg.App.KafkaUrl)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	} else {
		tx.Commit()
	}

	writeJSON(w, http.StatusOK, resp)

}

func (cfg *Config) balance(w http.ResponseWriter, r *http.Request) {
	var requestPayload m.JsonRequest
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		writeError(w, http.StatusNotAcceptable, errors.New("error unmarshaling request"))
		return
	}

	if requestPayload.Email == "" {
		writeError(w, http.StatusNotAcceptable, errors.New("email empty"))
		return
	}

	user, err := database.GetUserByEmail(cfg.Db, requestPayload.Email)
	if err != nil {
		writeError(w, http.StatusNotAcceptable, err)
		return
	}

	msg, err := cfg.Nts.Request("balance", []byte(fmt.Sprint(user.UserId)), 2*time.Second)
	if err != nil {
		if cfg.Nts.LastError() != nil {
			log.Printf("%v for request", cfg.Nts.LastError())
		}
		log.Printf("%v for request", err)
		writeError(w, http.StatusInternalServerError, errors.New("error trying to get balance"))
		return
	}

	log.Printf("Published [%s] : '%s'", "balance", requestPayload.Email)
	log.Printf("Received  [%v] : '%s'", msg.Subject, string(msg.Data))

	var userBalanceResponse = m.UserBalanceResponse{Email: requestPayload.Email}
	err = json.Unmarshal(msg.Data, &userBalanceResponse)
	if err != nil {
		writeError(w, http.StatusInternalServerError, errors.New("error trying to get balance"))
		return
	}

	writeJSON(w, http.StatusOK, userBalanceResponse)

}
