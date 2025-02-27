package httpserver

import (
	"errors"
	"net/http"
	"payment-gateway/internal/app/domain"
)

func (h HttpServer) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	var request TxStatusUpdate
	contentType := r.Header.Get("Content-Type")

	err := DecodeRequest(r, &request)

	if err != nil {
		BadRequest("unsupported-content-type", err, w)
		return
	}
	var status string
	var txId int64
	switch contentType {
	case "application/json":
		status = request.JsonStatus
		txId = request.JsonTxID
	case "text/xml", "application/xml":
		status = request.XmlStatus
		txId = request.XmlTxID
	default:
		BadRequest("unsupported-content-type", err, w)
		return

	}

	if errors.Is(ValidateStatus(status), domain.ErrInvalidTransactionStatus) {
		BadRequest("invalid-transaction-status", err, w)
		return
	}

	err = h.transactionPublisherService.PublishTransaction(r.Context(), txId, status, contentType)
	if err != nil {
		BadRequest("transaction-not-created", err, w)
		return
	}
	RespondOK(w, "async request for status update accepted", contentType, map[string]any{
		"transaction_id":  txId,
		"expected_status": status,
	})
}
