package models

import "github.com/google/uuid"

type Role string

const (
	ADMIN Role = "admin"
	GUEST      = "guest"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	Role     Role      `json:"role,omitempty"`
}

func NewUser(Username string, Password string, Role Role, ID uuid.UUID) *User {
	return &User{ID: ID, Username: Username, Password: Password, Role: Role}
}
