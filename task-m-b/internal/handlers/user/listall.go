package user

import (
	"log/slog"
	"net/http"
	"task-manager/internal/model"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type AllUserGetter interface {
	ListAllUsers() ([]model.User, error)
}

func All(log *slog.Logger, allUserGetter AllUserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.listall"
		log = log.With("operation", op, "request_id", middleware.GetReqID(r.Context()))

		users, err := allUserGetter.ListAllUsers()
		if err != nil {
			log.Error("Error listing users", "error", err)

			render.JSON(w, r, "internal server error")

			return
		}
		log.Info("got users")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]interface{}{
			"status": "OK",
			"data":   users,
		})
	}
}
