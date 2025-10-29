package user

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"task-manager/internal/lib/logger/sl"
	resp "task-manager/internal/lib/response"
	"task-manager/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type UserDeleter interface {
	DeleteUser(name string) error
}

// Delete handles deleting a user by email
func Delete(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Delete"

		log = log.With(
			"op", op,
			"request_id", middleware.GetReqID(r.Context()))

		raw := chi.URLParam(r, "user_name")
		email, errUnescape := url.PathUnescape(raw)
		if errUnescape != nil {
			email = raw
		}
		email = strings.TrimSpace(email)

		if email == "" {
			log.Error("name is empty")

			render.JSON(w, r, resp.Error("name email"))

			return
		}

		err := userDeleter.DeleteUser(email)
		if err != nil {
			if err == storage.ErrNotFound {
				log.Error("user not found", sl.Err(err))
				render.JSON(w, r, resp.Error("not found"))
				return
			}
			log.Error("failed to delete user", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		render.JSON(w, r, resp.OK())
	}
}
