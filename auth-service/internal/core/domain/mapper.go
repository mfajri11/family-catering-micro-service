package domain

import (
	"database/sql"
	"time"
)

// tags level technologies are allowed in mapper struct
// regex for validation password ^(?=.*[0-9])(?=.*[a-zA-Z])(=?.*[\!\@\#\$\%\^\&\*\(\)\_\+\[\]\{\}\|\\\:\;\"\'\<\>\,\.\?\/\`\~]).{8,}$
// TODO: add validation for password (min 8 char at least one numbers, one letter and one symbol must be appeared from given password)
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // must be a valid email
	Password string `json:"password" validate:"required"`    // min char length is 8 char and consist of uppercase, lowercase and symbol (must be in encrypt form)
}

type LogoutRequest struct {
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	// Confirm password is handled by front end
}

type SessionData struct {
	SID          string    `redis:"sid"`
	Username     string    `redis:"username"`
	Email        string    `redis:"email"`
	RefreshToken string    `redis:"refresh_token"`
	IsValid      bool      `redis:"is_valid"`
	ExpiredAt    time.Time `redis:"expired_at"`
	CreatedAt    time.Time `redis:"created_at"`
	UpdatedAt    time.Time `redis:"updated_at"`
}

type LoginResponse struct {
	SID           string
	AccessToken   string
	RefreshToken  string
	MultipleLogin bool
}

type OwnerResponse struct {
	Id             int64
	Name           string
	Email          string
	PhoneNumber    string
	DateOfBirth    sql.NullString
	HashedPassword string
}
