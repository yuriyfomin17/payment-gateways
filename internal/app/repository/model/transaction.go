package model

import (
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"payment-gateway/internal/app/domain"
	"time"
)

type TransactionType int

const (
	Withdraw TransactionType = iota
	Deposit
)

func (t TransactionType) String() string {
	switch t {
	case Withdraw:
		return "withdraw"
	case Deposit:
		return "deposit"
	default:
		return "unknown"
	}
}

type TransactionStatus int

const (
	Pending TransactionStatus = iota
	Completed
	Failed
)

func (s TransactionStatus) String() string {
	switch s {
	case Pending:
		return "pending"
	case Completed:
		return "completed"
	case Failed:
		return "failed"
	default:
		return "unknown"
	}
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	ID        int64           `bun:",pk,autoincrement"`
	Amount    decimal.Decimal `bun:",notnull"`
	Type      string          `bun:",notnull"`
	Status    string          `bun:",notnull"`
	CreatedAt time.Time       `bun:",default:current_timestamp"`
	GatewayID int             `bun:",notnull"`
	CountryID int             `bun:",notnull"`
	UserID    int             `bun:",notnull"`
}

func (t *Transaction) Validate() error {
	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return domain.ErrNegativeTransaction
	}

	return nil
}
