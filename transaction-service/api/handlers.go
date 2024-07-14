package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"transaction-service/database"
	"transaction-service/model"
	m "transaction-service/model"

	"github.com/shopspring/decimal"
)

var ErrClientUser = errors.New("user does not exist")
var ErrClientAmount = errors.New("amount error, check if it's a number with max two decimals")
var ErrServerDb = errors.New("there has been an error updating data")

func (cfg *Config) addMoney(w http.ResponseWriter, r *http.Request) {

	var requestPayload m.AddMoneyRequest
	// is request ok
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("request is not in good format"))
		return
	}

	// update balance, add transaction record, under db transaction
	// making sure balance is locked
	ctx := context.Background()
	tx, err := cfg.Db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		writeError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	newBalance, err := updateBalance(tx, cfg.Db, requestPayload.UserId, requestPayload.Amount)
	if err != nil {
		if errors.Is(err, ErrClientUser) || errors.Is(err, ErrClientAmount) || errors.Is(err, database.ErrDBalanceNegative) {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}
	tx.Commit()

	addMoneyResponse := model.AddMoneyResponse{
		UpdatedBalance: *newBalance,
	}

	writeJSON(w, http.StatusOK, addMoneyResponse)
}

func (cfg *Config) transferMoney(w http.ResponseWriter, r *http.Request) {

	var requestPayload m.TransferMoneyRequest
	// is request ok
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("request is not in good format"))
		return
	}

	// transfer money same as add money, it's just for multiple users instead of one
	transfer := []model.AddMoneyRequest{
		{
			UserId: requestPayload.FromUserId,
			Amount: requestPayload.Amount.Neg(),
		},
		{
			UserId: requestPayload.ToUserId,
			Amount: requestPayload.Amount,
		},
	}

	// update balances, add transaction records, under db transaction
	// making sure balance is locked
	ctx := context.Background()
	tx, err := cfg.Db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		writeError(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	for _, req := range transfer {
		_, err := updateBalance(tx, cfg.Db, req.UserId, requestPayload.Amount)
		if err != nil {
			if errors.Is(err, ErrClientUser) || errors.Is(err, ErrClientAmount) || errors.Is(err, database.ErrDBalanceNegative) {
				writeError(w, http.StatusBadRequest, err)
				return
			}

			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	tx.Commit()

	transferMoneyResponse := model.TransferMoneyResponse{}

	writeJSON(w, http.StatusOK, transferMoneyResponse)
}

// Updates balance, insert transaction and returns
// new balance for given userId
func updateBalance(tx *sql.Tx, db *sql.DB, userId int, amount decimal.Decimal) (*decimal.Decimal, error) {

	// we want max two decimals
	if amount != amount.Truncate(2) {
		return nil, ErrClientAmount
	}

	// does user exists
	_, err := database.GetBalance(db, userId)
	if err != nil {
		return nil, ErrClientUser
	}

	newAmount, err := database.TryUpdateBalanceWLock(tx, userId, amount)
	if err != nil {
		return nil, err
	}

	err = database.AddTransaction(tx, userId, amount)
	if err != nil {
		return nil, ErrServerDb
	}

	return newAmount, nil

}
