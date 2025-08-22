package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	"task-manager/internal/handlers/save/mocks"
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userSave := mocks.NewUserSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				userSave.
					On("SaveUser", tc.userName, tc.email).
					Return(int64(0), tc.mockError).
					Once()
			}
			log := slog.Default()
			handler := NewUser(log, userSave)

			inputs := fmt.Sprintf(`{"name":"%s", "email":"%s"}`, tc.userName, tc.email)

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
