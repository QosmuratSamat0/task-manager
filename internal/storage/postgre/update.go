package postgre

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) UpdateUserName(name string, newName string) (int64, error) {
	const op = "storage.postgre.UpdateUser"

	var id int64

	err := s.db.QueryRow(context.Background(),
		`UPDATE users SET user_name=$1, updated_at=now() WHERE user_name=$2`,
		name, newName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) UpdateTaskTitle(userId int64, title string) (int64, error) {
	const op = "storage.postgre.UpdateTaskTitle"

	var id int64
	err := s.db.QueryRow(context.Background(),
		`UPDATE tasks SET task_title=$1, updated_at=now() WHERE user_id=$2`, title, userId).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
