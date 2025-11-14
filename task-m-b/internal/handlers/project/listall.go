package project

import (
	"log/slog"
	"net/http"
	"task-manager/internal/model"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type AllProjectGetter interface {
	ListAllProject() ([]model.Project, error)
}

func All(log *slog.Logger, allProjectGetter AllProjectGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.project.listall"
		log = log.With("operation", op, "request_id", middleware.GetReqID(r.Context()))

		projects, err := allProjectGetter.ListAllProject()
		if err != nil {
			log.Error("Error listing project", "error", err)

			render.JSON(w, r, "internal server error")
		}
		log.Info("got projects")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   projects,
		})
	}
}
