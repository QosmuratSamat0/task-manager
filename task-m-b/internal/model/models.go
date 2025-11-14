package model

import "time"

type User struct {
    Id        int64      `json:"id"`
    UserName  string     `json:"user_name"`
    Email     string     `json:"email" validate:"required,email"`
    CreatedAt *time.Time `json:"created_at"`
    // Password contains raw password only for incoming requests; it is never returned.
    Password      string `json:"password,omitempty"`
    PasswordHash  string `json:"-"`
}

type Task struct {
	Id          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   *time.Time `json:"created_at"`
}

type Project struct {
	Id          int64      `json:"id"`
	OwnerId     int64      `json:"owner_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
