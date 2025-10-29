package task

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	resp "task-manager/internal/lib/response"
)

type TaskDeleter interface {
	DeleteTask(id int64) error
}

// Delete handles deleting a task by id
func Delete(log *slog.Logger, deleteTask TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.Delete"

		log = log.With("operation", op, "request_id", middleware.GetReqID(r.Context()))

		userIDstr := chi.URLParam(r, "user_id")

		userID, err := strconv.Atoi(userIDstr)

		if err != nil {
			log.Error("failed to parse user id")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		if userID < 0 {
			log.Error("user_id less than 0")

			render.JSON(w, r, resp.Error("invalid id"))

			return
		}

		err = deleteTask.DeleteTask(int64(userID))

		if err != nil {
			log.Error("failed to delete task")

			render.JSON(w, r, resp.Error("internal error"))

			return

		}

		render.JSON(w, r, resp.OK())
	}

}
