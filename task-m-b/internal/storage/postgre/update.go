package postgre

import (
    "context"
    "errors"
    "fmt"
    "task-manager/internal/model"
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

func (s *Storage) UpdateTaskTitle(title string, newTitle string) (int64, error) {
	const op = "storage.postgre.UpdateTaskTitle"

	var id int64
	err := s.db.QueryRow(context.Background(),
		`UPDATE tasks SET task_title=$1, updated_at=now() WHERE task_title=$2`, title, newTitle).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) UpdateTaskDescription(name string, description string) (int64, error) {
	const op = "storage.postgre.UpdateTaskDescription"

	var id int64
	err := s.db.QueryRow(context.Background(),
		`UPDATE project SET name=$1, updated_at=now() WHERE description=$2`, name, description,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
	}
	return id, nil
}

// UpdateTaskFields updates status and priority (string form) by task id
func (s *Storage) UpdateTaskFields(id int64, status string, priority string) (model.Task, error) {
    const op = "storage.postgre.UpdateTaskFields"

    // Reuse normalizePriority from save.go
    pri := normalizePriority(priority)

    var t model.Task
    err := s.db.QueryRow(context.Background(),
        `UPDATE tasks
           SET status = $1,
               priority = $2
         WHERE id = $3
         RETURNING id,
                   user_id,
                   title,
                   description,
                   status,
                   CASE priority WHEN 0 THEN 'low' WHEN 1 THEN 'medium' WHEN 2 THEN 'high' ELSE priority::text END AS priority,
                   deadline,
                   created_at`,
        status, pri, id,
    ).Scan(&t.Id, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Deadline, &t.CreatedAt)
    if err != nil {
        return model.Task{}, fmt.Errorf("%s: %w", op, err)
    }
    return t, nil
}
