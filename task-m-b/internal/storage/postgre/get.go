package postgre

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/model"
	"task-manager/internal/storage"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) User(name string) (model.User, error) {
    const op = "storage.postgre.User"

    var resUser model.User

    err := s.db.QueryRow(context.Background(),
        `SELECT id, user_name, email, created_at, password_hash FROM users WHERE user_name = $1`,
        name).Scan(&resUser.Id, &resUser.UserName, &resUser.Email, &resUser.CreatedAt, &resUser.PasswordHash)

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
        `SELECT id, user_id, title, description, status,
                CASE priority WHEN 0 THEN 'low' WHEN 1 THEN 'medium' WHEN 2 THEN 'high' ELSE priority::text END AS priority,
                deadline, created_at
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

func (s *Storage) Project(name string) (model.Project, error) {
	const op = "storage.postgre.Project"
	var p model.Project

	err := s.db.QueryRow(context.Background(),
		`SELECT id, owner_id, name, description, created_at, updated_at FROM project WHERE name = $1`, name,
	).Scan(
		&p.Id,
		&p.OwnerId,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Project{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}
	if err != nil {
		return model.Project{}, fmt.Errorf("%s: %w", op, err)
	}
	return p, nil
}

//func (s *Storage) Status(title string) []model.Status, error {
//	const op = "storage.postgre.Status"
//
//	var t model.Task
//
//	var tasks []model.Task
//
//	rows, err := s.db.Query(context.Background(),
//		`SELECT title, status FROM tasks WHERE title = $1`, title,
//	)
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//	defer rows.Close()
//	for rows.Next() {
//		if err := rows.Scan(&t.Title, &t.Status); err != nil {
//			return fmt.Errorf("%s: %w", op, err)
//		}
//	}
//	return rows.Err()
//}
