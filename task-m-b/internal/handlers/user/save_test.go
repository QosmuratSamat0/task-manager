package user

import (
	mock "github.com/stretchr/testify/mock"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	"task-manager/internal/handlers/user/mocks"
	"task-manager/internal/storage"
	"testing"
)

func TestSaveUser(t *testing.T) {
    cases := []struct {
        name      string
        userName  string
        email     string
        respError string
        mockError error
    }{
		{
			name:     "success",
			userName: "user1",
			email:    "user1@gmail.com",
		},
		{
			name:     "empty username",
			userName: "",
			email:    "user2@gmail.com",
		},
		{
			name:      "empty email",
			userName:  "user3",
			email:     "",
			respError: "email field is required",
		},
		{
			name:      "Invalid email",
			userName:  "user4",
			email:     "some invalid email.",
			respError: "invalid email format",
		},
		{
			name:      "SaveUser Error",
			userName:  "user5",
			email:     "user5@gmail.com",
			respError: "failed to save user",
			mockError: errors.New("unexpected error"),
		},
		{
			name:      "Duplicate email",
			userName:  "user6",
			email:     "user6@gmail.com",
			respError: "user already exists",
			mockError: storage.ErrExists,
		},
		{
			name:     "Uppercase email",
			userName: "user7",
			email:    "USER7@GMAIL.COM",
		},
		{
			name:     "Subdomain email",
			userName: "user8",
			email:    "user8@mail.example.com",
		},
		{
			name:     "Plus alias email",
			userName: "user9",
			email:    "user9+tag@gmail.com",
		},
		{
			name:     "Long email",
			userName: "user10",
			email:    "very.long.name.with.dots+tag@sub.example.co.uk",
		},
		{
			name:      "Invalid missing domain",
			userName:  "user11",
			email:     "user11@",
			respError: "invalid email format",
		},
		{
			name:      "Invalid missing local",
			userName:  "user12",
			email:     "@example.com",
			respError: "invalid email format",
		},
		{
			name:      "Invalid space in email",
			userName:  "user13",
			email:     "user 13@example.com",
			respError: "invalid email format",
		},
		{
			name:     "Unicode name",
			userName: "Самат",
			email:    "samat@example.com",
		},
		{
			name:      "SaveUser Error 2",
			userName:  "user15",
			email:     "user15@gmail.com",
			respError: "failed to save user",
			mockError: errors.New("db timeout"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userSave := mocks.NewUserSaver(t)

			if tc.respError == "" || tc.mockError != nil {
            userSave.
                On("SaveUser", tc.userName, tc.email, mock.AnythingOfType("string")).
                Return(int64(0), tc.mockError).
                Once()
        }
        log := slog.Default()
        handler := Save(log, userSave)

            inputs := fmt.Sprintf(`{"name":"%s", "email":"%s", "password":"password123"}`, tc.userName, tc.email)

			req, err := http.NewRequest(http.MethodPost, "/save", strings.NewReader(inputs))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp ResponseUser

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
