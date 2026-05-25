package postgres

import (
	"context"

	"github.com/google/uuid"
	"wallet-api/internal/domain"
)

type transactionRepo struct {
	db DBTX
}

func NewTransactionRepository(db DBTX) domain.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, idempotency_key, wallet_id, amount, operation_type, reference_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		tx.ID,
		tx.IdempotencyKey,
		tx.WalletID,
		tx.Amount,
		tx.OperationType,
		tx.ReferenceID,
		tx.CreatedAt,
	)

	if err != nil {
		// In a production environment, we should check if the error is a Postgres unique constraint violation
		// (SQLSTATE 23505) on idempotency_key, and return domain.ErrDuplicateTransaction
		// But for now, we return the generic error or you can use a library like 'pgx' error parsing.
		return err
	}

	return nil
}

func (r *transactionRepo) ListByWalletID(ctx context.Context, walletID uuid.UUID) ([]domain.Transaction, error) {
	query := `
		SELECT id, idempotency_key, wallet_id, amount, operation_type, reference_id, created_at
		FROM transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, walletID) // Assuming DBTX has QueryContext, wait let me add it to the interface
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.IdempotencyKey,
			&tx.WalletID,
			&tx.Amount,
			&tx.OperationType,
			&tx.ReferenceID,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}
