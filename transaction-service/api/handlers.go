package api

import (
	"context"
	"errors"
	"net/http"

	m "transaction-service/model"
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

}
