package api

import (
	"context"
	"errors"
	"net/http"

	"user-service/database"
	"user-service/messaging"
	m "user-service/model"
)

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

	ctx := context.Background()

	cfg.Tx, err = cfg.Db.BeginTx(ctx, nil)
	defer cfg.Tx.Rollback()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}

	resp, err := database.CreateUserWithTx(ctx, cfg.Tx, cfg.Db, requestPayload.Email)
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
