package domain

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	id              int64
	gatewayID       int
	countryID       int
	userID          int
	amount          decimal.Decimal
	transactionType string
	status          string
	dataFormat      string
	createdAt       time.Time
}

type TransactionData struct {
	ID              int64
	GatewayID       int
	CountryID       int
	UserID          int
	Currency        string
	Amount          decimal.Decimal
	TransactionType string
	Status          string
	DataFormat      string
	CreatedAt       time.Time
}

func NewTransaction(data TransactionData) (Transaction, error) {
	return Transaction{
		id:              data.ID,
		gatewayID:       data.GatewayID,
		countryID:       data.CountryID,
		userID:          data.UserID,
		amount:          data.Amount,
		transactionType: data.TransactionType,
		status:          data.Status,
		dataFormat:      data.DataFormat,
		createdAt:       data.CreatedAt,
	}, nil
}

func (t Transaction) ID() int64 {
	return t.id
}

func (t Transaction) Amount() decimal.Decimal {
	return t.amount
}

func (t Transaction) Type() string {
	return t.transactionType
}

func (t Transaction) Status() string {
	return t.status
}

func (t Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t Transaction) GatewayID() int {
	return t.gatewayID
}

func (t Transaction) CountryID() int {
	return t.countryID
}

func (t Transaction) UserID() int {
	return t.userID
}

func (t Transaction) DataFormat() string {
	return t.dataFormat
}

func (t Transaction) ToJSON() ([]byte, error) {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction to JSON: %w", err)
	}
	return jsonData, nil
}

// ToSOAP converts the Transaction to SOAP format.
func (t Transaction) ToSOAP() ([]byte, error) {

	// Create a SOAP envelope struct (adapt as needed for your SOAP schema)
	type SoapEnvelope struct {
		XMLName xml.Name `xml:"soap:Envelope"`
		Soap    string   `xml:"xmlns:soap,attr"`
		Body    struct {
			TransactionData Transaction `xml:"TransactionData"`
		}
	}

	soapEnvelope := SoapEnvelope{
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: struct {
			TransactionData Transaction `xml:"TransactionData"`
		}{TransactionData: t},
	}

	soapData, err := xml.MarshalIndent(soapEnvelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction to SOAP: %w", err)
	}
	return soapData, nil
}
