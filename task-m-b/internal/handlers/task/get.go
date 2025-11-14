package task

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	resp "task-manager/internal/lib/response"
	"task-manager/internal/model"
	"task-manager/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=TaskGetter
type TaskGetter interface {
	Task(id int) (model.Task, error)
}

// Get handles retrieving a task by id
func Get(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.Get"

		log = log.With("op", op, "request_id", middleware.GetReqID(r.Context()))

		userIDstr := chi.URLParam(r, "user_id")
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			log.Error("Error converting user_id to int", "user_id", userIDstr)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid user_id"))

			return
		}
		if userID < 0 {
			log.Error("user_id less than 0")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("user id less than 0"))

			return
		}

		task, err := taskGetter.Task(userID)
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("task not found")

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("task not found"))

			return
		}

		if err != nil {
			log.Error("failed to get task", "user_id", userID)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got task", "user_id", userID)

		render.Status(r, http.StatusOK)

		render.JSON(w, r, map[string]interface{}{
			"status": "ok",
			"task":   task,
		})

	}
}
