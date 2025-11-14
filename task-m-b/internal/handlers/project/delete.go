package project

import (
	"log/slog"
	"net/http"
	"strconv"
	resp "task-manager/internal/lib/response"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type ProjectDeleter interface {
	DeleteProject(ownerId int64) error
}

func Delete(log *slog.Logger, projectDeleter ProjectDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.project.Delete"

		log = log.With("operation", op,
			"request_id", middleware.GetReqID(r.Context()),
		)
		ownerId := chi.URLParam(r, "owner_id")
		id, err := strconv.Atoi(ownerId)
		if err != nil {
			log.Error("failed to parse owner id")

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "internal server error")

			return
		}

		if id <= 0 {
			log.Error("owner id must be positive")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "owner id must be positive")

			return
		}

		err = projectDeleter.DeleteProject(int64(id))

		if err != nil {
			log.Error("failed to delete project")

			render.JSON(w, r, "internal server error")

			return
		}

		render.JSON(w, r, resp.OK())

	}
}
