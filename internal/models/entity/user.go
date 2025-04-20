package entity

import (
	"database/sql/driver"
	"errors"
	"github.com/AnikinSimon/avito-test-backend/internal/models/dto/response"

	"github.com/google/uuid"
)

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

// for fast-checking.
var Roles = map[Role]bool{
	RoleEmployee:  true,
	RoleModerator: true,
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *Role) Scan(value interface{}) error {
	*r = Role(string(value.([]byte)))
	return nil
}

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     Role
}

func (u *User) ToResponse() *response.User {
	return &response.User{
		ID:    u.ID,
		Email: u.Email,
		Role:  string(u.Role),
	}
}

func (u *User) MarshalJSON() ([]byte, error) {
	return nil, errors.New("entity.User: direct JSON serialization forbidden, use response.User")
}
