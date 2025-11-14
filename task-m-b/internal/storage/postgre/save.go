package postgre

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"task-manager/internal/model"
	"task-manager/internal/storage"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) SaveUser(username string, email string, passwordHash string) (int64, error) {
    const op = "storage.postgre.SaveUser"

    var id int64
    err := s.db.QueryRow(context.Background(), "INSERT INTO users (user_name, email, password_hash) VALUES ($1, $2, $3) RETURNING id", username, email, passwordHash).Scan(&id)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
        }
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return id, nil
}

func (s *Storage) SaveTask(userID int64, title, description, status, priority string, deadline time.Time) (int64, error) {
    const op = "storage.postgre.SaveTask"

    var id int64

    // Normalize priority to integer
    pri := normalizePriority(priority)

    var deadlineParam any
    if deadline.IsZero() {
        deadlineParam = nil
    } else {
        deadlineParam = deadline
    }

    err := s.db.QueryRow(context.Background(),
        `
            INSERT INTO tasks (user_id, title, description, status, priority, deadline)
            VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
        `, userID, title, description, status, pri, deadlineParam).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) SaveProject(p model.Project) (int64, error) {
	const op = "storage.postgre.SaveProject"
	var id int64
	err := s.db.QueryRow(context.Background(),
		`INSERT INTO project (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id`,
		p.OwnerId, p.Name, p.Description,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func normalizePriority(p string) int {
	ps := strings.TrimSpace(strings.ToLower(p))
	switch ps {
	case "low":
		return 0
	case "normal", "medium":
		return 1
	case "high":
		return 2
	}
	if n, err := strconv.Atoi(ps); err == nil {
		return n
	}
	return 0
}
