package pgrepo

import (
	"context"
	"fmt"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
	"payment-gateway/internal/pkg/pg"
)

// TransactionRepo manages transaction-related DB operations.
type TransactionRepo struct {
	db *pg.DB
}

// NewTransactionRepo constructs a new TransactionRepo with the given database connection.
func NewTransactionRepo(db *pg.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}
func (repo TransactionRepo) UpdateTransactionStatus(ctx context.Context, transactionID int64, status string) error {
	res, err := repo.db.NewUpdate().
		Model(&model.Transaction{}).
		Set("status = ?", status).
		Where("id = ?", transactionID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	interactedRows, err := res.RowsAffected()
	if err != nil || interactedRows == 0 {
		return domain.ErrTransactionNotFound
	}

	return nil
}
func (repo TransactionRepo) FetchTxId(ctx context.Context) (int64, error) {
	var txId int64
	err := repo.db.NewSelect().ColumnExpr("nextval('transactions_tx_counter_seq')").Scan(ctx, &txId)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch next transaction ID from sequence: %w", err)
	}
	return txId, nil
}

func (repo TransactionRepo) CreateTransaction(ctx context.Context, transaction domain.Transaction) error {
	modelTransaction, err := domainToTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to convert transaction to model: %w", err)
	}
	_, err = repo.db.NewInsert().Model(&modelTransaction).
		On("CONFLICT (id) DO UPDATE"). // Handle potential ID conflicts
		Set("gateway_id = EXCLUDED.gateway_id, country_id = EXCLUDED.country_id, user_id = EXCLUDED.user_id, amount = EXCLUDED.amount, type = EXCLUDED.type, status = EXCLUDED.status, created_at = EXCLUDED.created_at").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}
