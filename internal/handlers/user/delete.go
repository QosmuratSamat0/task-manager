package user

import (
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "log/slog"
    "net/http"
    resp "task-manager/internal/lib/response"
)

type UserDeleter interface {
    DeleteUser(email string) error
}

// Delete handles deleting a user by email
func Delete(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.user.Delete"

        log = log.With(
            "op", op,
            "request_id", middleware.GetReqID(r.Context()))

        email := chi.URLParam(r, "email")

        if email == "" {
            log.Error("email is empty")

            render.JSON(w, r, resp.Error("empty email"))

            return
        }

        err := userDeleter.DeleteUser(email)
        if err != nil {
            log.Error("failed to delete user")

            render.JSON(w, r, resp.Error("internal error"))

            return
        }

        render.JSON(w, r, resp.OK())
    }
}

