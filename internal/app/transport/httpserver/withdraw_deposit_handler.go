package httpserver

import (
	"errors"
	"net/http"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
)

func (h HttpServer) DepositHandler(w http.ResponseWriter, r *http.Request) {
	h.handleTransaction(w, r, model.Deposit.String())
}

func (h HttpServer) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	h.handleTransaction(w, r, model.Withdraw.String())
}

func (h HttpServer) handleTransaction(w http.ResponseWriter, r *http.Request, txType string) {
	var request TransactionRequest

	err := DecodeRequest(r, &request)
	if err != nil {
		BadRequest("unsupported-content-type", err, w)
		return
	}
	txId, err := h.userService.FetchTxId(r.Context())
	if err != nil && errors.Is(err, domain.ErrTransactionNotCreated) {
		BadRequest("transaction-not-created", err, w)
		return
	}

	if err != nil {
		InternalError("transaction-not-created", err, w)
		return
	}
	// TODO:
	// implement support for SOAP
	err = h.userService.ExecuteTransaction(r.Context(), domain.TransactionData{
		ID:              txId,
		UserID:          request.UserID,
		Amount:          request.Amount,
		Currency:        request.Currency,
		DataFormat:      "JSON",
		TransactionType: txType,
	})
	if err != nil {
		BadRequest("transaction-not-created", err, w)
		return
	}
	RespondOK(w, "success", "application/json", map[string]any{
		"transaction_id":     txId,
		"transaction_status": "pending",
	})
}
