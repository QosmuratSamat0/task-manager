package task

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"task-manager/internal/lib/logger/sl"
	resp "task-manager/internal/lib/response"
	"task-manager/internal/storage"
	"time"
)

type RequestTask struct {
	UserID      int64     `json:"user_id" validate:"required"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
}

type ResponseTask struct {
	resp.Response
	UserID int64 `json:"user_id"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=TaskSaver
type TaskSaver interface {
	SaveTask(userID int64, title, description, status, priority string, deadline time.Time) (int64, error)
}

// Save handles creating a task
func Save(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.Save"

		log = log.With(
			"op", op,
			"request_id", middleware.GetReqID(r.Context()))

		var req RequestTask

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.UserID == 0 {
			log.Error("missing user_id", sl.Err(err))

			render.JSON(w, r, resp.Error("missing user_id"))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

        id, err := taskSaver.SaveTask(req.UserID, req.Title, req.Description, req.Status, req.Priority, req.Deadline)
        if err != nil {
            if errors.Is(err, storage.ErrExists) {
                log.Error("duplicate task title for user", sl.Err(err))

                render.JSON(w, r, resp.Error("task with this title already exists for this user"))

                return
            }
            log.Error("failed to save task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save task"))

			return
		}

		slog.Info("saved task", slog.Int64("user_id", id))

		render.JSON(w, r, ResponseTask{
			Response: resp.OK(),
			UserID:   req.UserID,
		})

	}
}
