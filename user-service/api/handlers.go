package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"user-service/database"
	"user-service/messaging"
	m "user-service/model"
)

func (cfg *Config) getTime(w http.ResponseWriter, r *http.Request) {

	time, err := json.Marshal(time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp := m.JsonResponse{
		Error:   false,
		Message: "",
		Data:    string(time),
	}

	writeJSON(w, http.StatusOK, resp)

}

func (cfg *Config) createUser(w http.ResponseWriter, r *http.Request) {

	var requestPayload m.JsonRequest
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		writeError(w, http.StatusNotAcceptable, errors.New("error unmarshaling Url"))
		return
	}

	if requestPayload.Email == "" {
		writeError(w, http.StatusNotAcceptable, errors.New("email empty"))
		return
	}

	cfg.Tx, err = cfg.Db.BeginTx(context.Background(), nil)
	defer cfg.Tx.Rollback()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp, err := database.CreateUserWithTx(context.Background(), cfg.Tx, cfg.Db, requestPayload.Email)
	if err != nil {
		if errors.Is(err, database.ErrDuplicate) {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// create in kafka
	err = messaging.PublishJSON("user-created", resp, "localhost:9092")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	} else {
		cfg.Tx.Commit()
	}

	writeJSON(w, http.StatusOK, resp)

}
