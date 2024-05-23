package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"test-task/internal/entity"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) RateCreate(ctx context.Context, timestamp int64, ask, bid entity.Rate) error {
	q := `INSERT INTO "rates" ("id", "timestamp", "ask", "bid") VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, q, uuid.New(), timestamp, ask, bid)
	if err != nil {
		return fmt.Errorf("postgresql exec: %w", err)
	}
	return nil
}
