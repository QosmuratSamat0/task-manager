package postgre

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"task-manager/internal/storage"
	"time"
)

func (s *Storage) SaveUser(name string, email string) (int64, error) {
	const op = "storage.postgre.SaveUser"

	var id int64
	err := s.db.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", name, email).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SaveTask(userID int, title, description, status, priority string, deadline time.Time) (int64, error) {
	const op = "storage.postgre.SaveTask"

	var id int64

	err := s.db.QueryRow(context.Background(),
		`
			INSERT INTO tasks (user_id, title, description, status, priority, deadline) VALUES ($1, $2, $3, $4, $5)
		`, userID, title, description, status, priority, deadline).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
