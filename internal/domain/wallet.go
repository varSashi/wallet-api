package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Wallet stores the user's cached balance reference.
type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   int64     `json:"balance"` // In cents to avoid floating-point errors
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WalletRepository dictates how we access the wallets table.
type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*Wallet, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, newBalance int64) error
}
