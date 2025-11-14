package postgre

import (
	"context"
	"fmt"
	"task-manager/internal/storage"
)

func (s *Storage) DeleteUser(name string) error {
	const op = "storage.postgre.DeleteUser"

	res, err := s.db.Exec(context.Background(), `DELETE FROM users WHERE user_name = $1`, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) DeleteTask(id int64) error {
	const op = "storage.postgre.DeleteTask"

	res, err := s.db.Exec(context.Background(), `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) DeleteProject(ownerId int64) error {
	const op = "storage.postgre.DeleteProject"

	res, err := s.db.Exec(context.Background(), `DELETE FROM project WHERE owner_id = $1`, ownerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return nil
}
