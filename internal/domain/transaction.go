package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationTypeDeposit     OperationType = "DEPOSIT"
	OperationTypeTransferOut OperationType = "TRANSFER_OUT"
	OperationTypeTransferIn  OperationType = "TRANSFER_IN"
)

// Transaction represents our Ledger. An immutable entry in the statement.
type Transaction struct {
	ID             uuid.UUID      `json:"id"`
	IdempotencyKey string         `json:"idempotency_key"`
	WalletID       uuid.UUID      `json:"wallet_id"`
	Amount         int64          `json:"amount"` // Negative values indicate outgoing funds
	OperationType  OperationType  `json:"operation_type"`
	ReferenceID    *uuid.UUID     `json:"reference_id,omitempty"` // Links sender and receiver transactions
	CreatedAt      time.Time      `json:"created_at"`
}

// Domain Errors. All layers will use these standardized errors.
var (
	ErrInsufficientFunds = errors.New("insufficient funds for the transaction")
	ErrDuplicateTransaction = errors.New("duplicate transaction (idempotency)")
	ErrWalletNotFound = errors.New("wallet not found")
	ErrUserNotFound = errors.New("user not found")
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	ListByWalletID(ctx context.Context, walletID uuid.UUID) ([]Transaction, error)
}
