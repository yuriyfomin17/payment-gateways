package httpserver

import (
	"encoding/xml"
	"github.com/shopspring/decimal"
	"payment-gateway/internal/app/domain"
)

type TransactionRequest struct {
	Amount   decimal.Decimal `json:"amount"`
	UserID   int             `json:"user_id"`
	Currency string          `json:"currency"`
}

type TxStatusUpdate struct {
	JsonTxID   int64  `json:"tx_id"`
	JsonStatus string `json:"status"`

	XMLName   xml.Name `xml:"tx_status"`
	XmlTxID   int64    `xml:"tx_id"`
	XmlStatus string   `xml:"status"`
}

type GatewayPriorityUpdate struct {
	GtId     int64  `json:"gt_id"`
	Priority string `json:"priority"`
}

func ValidateStatus(status string) error {
	switch status {
	case "pending", "completed", "failed":
		return nil
	}
	return domain.ErrInvalidTransactionStatus
}
func ValidateGatewayPriority(priority string) error {
	switch priority {
	case "high", "medium", "low":
		return nil
	}
	return domain.ErrInvalidGatewayPriority
}

// APIResponse a standard response structure for the APIs
type APIResponse struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Message    string `json:"message" xml:"message"`
	Data       any    `json:"data,omitempty" xml:"data,omitempty"`
}
