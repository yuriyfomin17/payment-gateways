package domain

import "errors"

var (
	ErrNegativeTransaction              = errors.New("amount must be positive")
	ErrInvalidTransactionType           = errors.New("invalid transaction type")
	ErrInvalidTransactionStatus         = errors.New("invalid transaction status")
	ErrInvalidGatewayPriority           = errors.New("invalid gateway priority")
	ErrNotFound                         = errors.New("not found")
	ErrUserNotFound                     = errors.New("user not found")
	ErrTransactionNotCreated            = errors.New("transaction not created")
	ErrTransactionStatus                = errors.New("transaction status update not executed")
	ErrTransactionNotPublished          = errors.New("transaction not published")
	ErrTransactionNotFound              = errors.New("transaction not found")
	ErrRedisTransactionNotCreated       = errors.New("redis transaction not created")
	ErrRedisCouldNotExtractTransactions = errors.New("could not extract transactions from redis")
	ErrRedisCreateTransaction           = errors.New("could not create transaction in redis")
)
