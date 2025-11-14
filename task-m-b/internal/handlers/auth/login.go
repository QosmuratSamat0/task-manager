package auth

import (
    "log/slog"
    "net/http"
    resp "task-manager/internal/lib/response"
    "task-manager/internal/model"

    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "golang.org/x/crypto/bcrypt"
)

// UserGetter provides user retrieval by username including password hash.
type UserGetter interface {
    User(name string) (model.User, error)
}

type loginRequest struct {
    UserName string `json:"user_name"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

func Login(log *slog.Logger, getter UserGetter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.auth.Login"
        log = log.With("op", op, "request_id", middleware.GetReqID(r.Context()))

        var req loginRequest
        if err := render.DecodeJSON(r.Body, &req); err != nil {
            render.JSON(w, r, resp.Error("failed to decode request"))
            return
        }

        username := req.UserName
        if username == "" {
            username = req.Name
        }
        if username == "" || req.Password == "" {
            render.Status(r, http.StatusBadRequest)
            render.JSON(w, r, resp.Error("missing credentials"))
            return
        }

        user, err := getter.User(username)
        if err != nil {
            render.Status(r, http.StatusUnauthorized)
            render.JSON(w, r, resp.Error("invalid credentials"))
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
            render.Status(r, http.StatusUnauthorized)
            render.JSON(w, r, resp.Error("invalid credentials"))
            return
        }

        // Return sanitized user
        render.Status(r, http.StatusOK)
        render.JSON(w, r, map[string]interface{}{
            "status": "OK",
            "data": map[string]interface{}{
                "id":        user.Id,
                "user_name": user.UserName,
                "email":     user.Email,
                "created_at": user.CreatedAt,
            },
        })
    }
}

