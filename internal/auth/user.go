package auth

import "github.com/google/uuid"

// Role model info
// @Description Value that represents user's right in the service.
type Role string

const (
	ADMIN Role = "admin" // This role provides full access to the API
	GUEST Role = "guest" // This role provides restricted access to the API - client gets only methods for observation.
)

// User model info
// @Description Model that represents user's model, also the content of a body in the request to be signed up.
type User struct {
	ID       uuid.UUID `json:"id,omitempty"`       // id of the user that's server-side generated
	Username string    `json:"username,omitempty"` // username of the client. Must be unique
	Password string    `json:"password,omitempty"` // password of the client that is hashed on the server side
	Role     Role      `json:"role,omitempty"`     // role of the user. Either guest or admin.
}

func NewUser(Username string, Password string, Role Role, ID uuid.UUID) *User {
	return &User{ID: ID, Username: Username, Password: Password, Role: Role}
}
