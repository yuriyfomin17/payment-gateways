//go:generate mockery

package httpserver

import (
	"context"
	"payment-gateway/internal/app/domain"
)

type UserService interface {
	ExecuteTransaction(ctx context.Context, tx domain.TransactionData) error
	FetchTxId(ctx context.Context) (int64, error)
}

type GatewayService interface {
	UpdateGatewayPriority(ctx context.Context, gatewayID int64, priority string) error
}
