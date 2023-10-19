package domain

import "time"

type Session struct {
	SID          string    `redis:"sid"`
	Email        string    `redis:"email"`
	UserID       string    `redis:"user_id"`
	RefreshToken string    `redis:"refresh_token"`
	IsValid      bool      `redis:"is_valid"`
	ExpiredAt    time.Time `redis:"expired_at"`
	CreatedAt    time.Time `redis:"created_at"`
	UpdatedAt    time.Time `redis:"updated_at"`
}
