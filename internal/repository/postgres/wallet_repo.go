package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"wallet-api/internal/domain"
)

// DBTX is an interface that matches both *sql.DB and *sql.Tx.
// This allows our repository to run queries either alone or inside an ACID transaction.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type walletRepo struct {
	db DBTX
}

func NewWalletRepository(db DBTX) domain.WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, wallet.ID, wallet.UserID, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt)
	return err
}

func (r *walletRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE id = $1`
	return r.scanWallet(r.db.QueryRowContext(ctx, query, id))
}

func (r *walletRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = $1`
	return r.scanWallet(r.db.QueryRowContext(ctx, query, userID))
}

// The "FOR UPDATE" locks the row in PostgreSQL
func (r *walletRepo) GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	query := `SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE id = $1 FOR UPDATE`
	return r.scanWallet(r.db.QueryRowContext(ctx, query, id))
}

func (r *walletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, newBalance int64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newBalance, id)
	return err
}

// Helper function to avoid repeating the Scan logic
func (r *walletRepo) scanWallet(row *sql.Row) (*domain.Wallet, error) {
	var w domain.Wallet
	err := row.Scan(&w.ID, &w.UserID, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}
	return &w, nil
}
