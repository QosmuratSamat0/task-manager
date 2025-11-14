package postgre

import (
	"context"
	"fmt"
	"task-manager/internal/model"
)

func (s *Storage) ListAllTasks() ([]model.Task, error) {
    const ap = "storage.postgres.ListAllTasks"

    rows, err := s.db.Query(context.Background(),
        `SELECT id,
                user_id,
                title,
                description,
                status,
                CASE priority WHEN 0 THEN 'low' WHEN 1 THEN 'medium' WHEN 2 THEN 'high' ELSE priority::text END AS priority,
                deadline,
                created_at
         FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ap, err)
	}

	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		err = rows.Scan(&t.Id, &t.UserID, &t.Title, &t.Description,
			&t.Status, &t.Priority, &t.Deadline, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ap, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", ap, err)
	}
	return tasks, nil
}

// ListTasksByUser returns tasks owned by a specific user id
func (s *Storage) ListTasksByUser(userID int64) ([]model.Task, error) {
    const ap = "storage.postgres.ListTasksByUser"

    rows, err := s.db.Query(context.Background(),
        `SELECT id,
                user_id,
                title,
                description,
                status,
                CASE priority WHEN 0 THEN 'low' WHEN 1 THEN 'medium' WHEN 2 THEN 'high' ELSE priority::text END AS priority,
                deadline,
                created_at
         FROM tasks
         WHERE user_id = $1`, userID)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", ap, err)
    }
    defer rows.Close()

    var tasks []model.Task
    for rows.Next() {
        var t model.Task
        if err := rows.Scan(&t.Id, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Deadline, &t.CreatedAt); err != nil {
            return nil, fmt.Errorf("%s: %w", ap, err)
        }
        tasks = append(tasks, t)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", ap, err)
    }
    return tasks, nil
}

func (s *Storage) ListAllUsers() ([]model.User, error) {
	const ap = "storage.postgres.ListAllUsers"

	db, err := s.db.Query(context.Background(), `SELECT id, user_name, email, created_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ap, err)
	}
	defer db.Close()
	var user []model.User
	for db.Next() {
		var u model.User
		err = db.Scan(&u.Id, &u.UserName, &u.Email, &u.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ap, err)
		}
		user = append(user, u)
	}
	if err := db.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", ap, err)
	}
	return user, nil
}

func (s *Storage) ListAllProjects() ([]model.Project, error) {
	const op = "storage.postgres.ListAllProjects"

	rows, err := s.db.Query(context.Background(),
		`SELECT id, owner_id, name, description, created_at, updated_at FROM project`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.Id,
			&p.OwnerId,
			&p.Name,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return projects, nil
}
