//go:generate mockery

package services

import (
	"context"
	"payment-gateway/internal/app/domain"
)

type DataEncryptor interface {
	MaskData(data []byte) ([]byte, error)
	UnmaskData(data []byte) ([]byte, error)
}

type FaultTolerance interface {
	PublishWithCircuitBreaker(operation func() error) error
	RetryOperation(operation func() error, maxRetries int) error
}

type RabbitMqService interface {
	PublishJsonData(jsonData []byte) error
	PublishSoapData(soapData []byte) error
	GetJsonMessage() chan []byte
	GetSoapMessage() chan []byte
}

type TransactionPublisherService interface {
	PublishTransaction(ctx context.Context, txId int64, statusToUpdate, dataFormat string) error
}

type RedisService interface {
	CreateFailedCallbackTransaction(ctx context.Context, txId int64, statusToUpdate, dataFormat string) error
	GetListOfFailedCallbackTransactionsToProcess(ctx context.Context) ([]domain.TransactionData, error)
	DeleteFailedCallbackTransaction(ctx context.Context, txId int64) error
}

type CountryRepository interface {
	GetCountryByID(ctx context.Context, countryID int, currencyId string) (domain.Country, error)
}

type TransactionRepository interface {
	FetchTxId(ctx context.Context) (int64, error)
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status string) error
	CreateTransaction(ctx context.Context, transaction domain.Transaction) error
}

type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (domain.User, error)
}

type GatewayRepository interface {
	UpdateGatewayPriority(ctx context.Context, gatewayID int64, priority string) error
	GetSupportedGatewaysByCountrySortedByPriorities(ctx context.Context, countryId int, dataFormat string) ([]domain.Gateway, error)
}
