package postgre

import (
    "context"
    "errors"
    "fmt"
    "github.com/jackc/pgx/v5"
    "task-manager/internal/model"
    "task-manager/internal/storage"
)

func (s *Storage) User(email string) (model.User, error) {
	const op = "storage.postgre.User"

	var resUser model.User

	err := s.db.QueryRow(context.Background(),
		`
			SELECT id, user_name, email, created_at FROM users WHERE email = $1
			`, email).Scan(&resUser.Id, &resUser.UserName, &resUser.Email, &resUser.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return resUser, nil
}

func (s *Storage) Task(id int) (model.Task, error) {
	const op = "storage.postgre.task"

	var task model.Task

	err := s.db.QueryRow(context.Background(),
		`SELECT id, user_id, title, description, status, priority, deadline, created_at
         FROM tasks
         WHERE id = $1`,
		id,
	).Scan(
		&task.Id,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.Deadline,
		&task.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Task{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}
	if err != nil {
		return model.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}
