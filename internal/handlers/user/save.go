package user

import (
    "errors"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "github.com/go-playground/validator/v10"
    "log/slog"
    "net/http"
    "net/mail"
    "strings"
    "task-manager/internal/lib/logger/sl"
    resp "task-manager/internal/lib/response"
    "task-manager/internal/storage"
)

type RequestUser struct {
    UserName string `json:"name"`
    Email    string `json:"email" validate:"required,email"`
}

type ResponseUser struct {
    resp.Response
    Email string `json:"email"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserSaver
type UserSaver interface {
    SaveUser(name string, email string) (int64, error)
}

// Save handles creating a user
func Save(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.user.Save"

        log = log.With(
            "op", op,
            "request_id", middleware.GetReqID(r.Context()),
        )

        var req RequestUser

        err := render.DecodeJSON(r.Body, &req)
        if err != nil {
            log.Error("failed to decode request body", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to decode request"))
            return
        }

        log.Info("request body decoded", slog.Any("request", req))

        if req.Email == "" {
            log.Error("email field is required")

            render.JSON(w, r, resp.Error("email field is required"))
            return
        }
        req.Email = strings.TrimSpace(req.Email)

        if _, err := mail.ParseAddress(req.Email); err != nil {
            log.Info("invalid email format", sl.Err(err))

            render.JSON(w, r, resp.Error("invalid email format"))
            return
        }

        if err := validator.New().Struct(req); err != nil {

            validateErr := err.(validator.ValidationErrors)
            log.Error("failed to validate request", sl.Err(err))

            render.JSON(w, r, resp.ValidationError(validateErr))
            return
        }

        id, err := userSaver.SaveUser(req.UserName, req.Email)
        if err != nil {
            if errors.Is(err, storage.ErrExists) {
                log.Error("user already exists", sl.Err(err))

                render.JSON(w, r, resp.Error("user already exists"))

                return
            }
            log.Error("failed to save user", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to save user"))

            return
        }

        log.Info("user saved", slog.Int64("id", id))

        render.JSON(w, r, ResponseUser{
            Response: resp.OK(),
            Email:    req.Email,
        })
    }
}

