package postgre

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgre.New"

	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

    return &Storage{db: db}, nil
}

// Ping checks database connectivity
func (s *Storage) Ping(ctx context.Context) error {
    return s.db.Ping(ctx)
}
