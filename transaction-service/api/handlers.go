package api

import (
	"context"
	"errors"
	"net/http"
	"transaction-service/database"
	"transaction-service/model"
	m "transaction-service/model"
)

var UserErrorNotExists = errors.New("user already exists")

func (cfg *Config) addMoney(w http.ResponseWriter, r *http.Request) {

	var requestPayload m.AddMoneyRequest
	// is request ok
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("request is not in good format"))
		return
	}

	// is amount ok
	amount := requestPayload.Amount

	// we want max two decimals
	if amount != amount.Truncate(2) {
		writeError(w, http.StatusBadRequest, errors.New("amount cannot have more than two decimal digits"))
		return
	}

	// does user exists
	_, err = database.GetBalance(cfg.Db, requestPayload.UserId)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("user not found"))
		return
	}

	// start locking things, we need to update balance and add transaction
	ctx := context.Background()
	tx, err := cfg.Db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		writeError(w, http.StatusInternalServerError, errors.New("internal database error"))
		return
	}

	newAmount, err := database.TryUpdateBalanceWLock(tx, requestPayload.UserId, amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = database.AddTransaction(tx, requestPayload.UserId, amount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	addMoneyResponse := model.AddMoneyResponse{
		UpdatedBalance: *newAmount,
	}
	tx.Commit()

	writeJSON(w, http.StatusOK, addMoneyResponse)

}

func (cfg *Config) transferMoney(w http.ResponseWriter, r *http.Request) {
}
