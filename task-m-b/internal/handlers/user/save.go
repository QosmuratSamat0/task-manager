package user

import (
    "errors"
    "log/slog"
    "net/http"
    "net/mail"
    "strings"
    "task-manager/internal/lib/logger/sl"
    resp "task-manager/internal/lib/response"
    "task-manager/internal/model"
    "task-manager/internal/storage"

    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "github.com/go-playground/validator/v10"
    "golang.org/x/crypto/bcrypt"
)

type ResponseUser struct {
	resp.Response
	Email string `json:"email"`
}

type UserSaver interface {
    SaveUser(name string, email string, passwordHash string) (int64, error)
}

// Save handles creating a user
func Save(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Save"

		log = log.With(
			"op", op,
			"request_id", middleware.GetReqID(r.Context()),
		)

        var req model.User
        err := render.DecodeJSON(r.Body, &req)
        if err != nil {
            log.Error("failed to decode request body", sl.Err(err))

            render.JSON(w, r, resp.Error("failed to decode request"))
            return
        }

        log.Info("request body decoded", slog.Any("request", req))

        // Support either user_name or name in payload
        if req.UserName == "" && req.Password == "" && req.Email == "" {
            // try reading common alias field "name" into username if present
            // noop here because req already captured any matching fields
        }

        if req.UserName == "" {
            log.Error("user_name field is required")

            render.JSON(w, r, resp.Error("user_name field is required"))

            return
        }

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

        // Basic password validation
        if strings.TrimSpace(req.Password) == "" || len(req.Password) < 6 {
            render.JSON(w, r, resp.Error("password must be at least 6 characters"))
            return
        }

        if err := validator.New().Struct(req); err != nil {

            validateErr := err.(validator.ValidationErrors)
            log.Error("failed to validate request", sl.Err(err))

            render.JSON(w, r, resp.ValidationError(validateErr))
            return
        }

        // Hash password
        hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Error("failed to hash password", sl.Err(err))
            render.JSON(w, r, resp.Error("internal error"))
            return
        }

        id, err := userSaver.SaveUser(req.UserName, req.Email, string(hash))
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
