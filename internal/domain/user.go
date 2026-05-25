package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the base entity of an application user.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the contract that the data layer (Postgres) must implement.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}
