package task

import (
	"log/slog"
	"net/http"
	"task-manager/internal/model"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type AllTaskGetter interface {
	ListAllTasks() ([]model.Task, error)
}

func All(log *slog.Logger, allTaskGetter AllTaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.listall"
		log = log.With("operation", op, "request_id", middleware.GetReqID(r.Context()))

		tasks, err := allTaskGetter.ListAllTasks()
		if err != nil {
			log.Error("Error listing tasks", "error", err)

			render.JSON(w, r, "internal server error")
		}
		log.Info("got tasks")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   tasks,
		})
	}
}
