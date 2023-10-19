package ports

import (
	"context"

	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/domain"
)

type IAuthService interface {
	Login(ctx context.Context, req domain.LoginRequest) (resp *domain.LoginResponse, err error)
	// TODO: Logout
	// TODO: Forgot password (when user already Logged in)
	// TODO: Reset password
	// TODO: Renew access token
	// TODO: Generate OTP
	// TODO: validate OTP
	// TODO: Send OTP (by email)
	// TODO: Password less login (magic link)
	// TODO: Revoke session
	// TODO: Create License
	// TODO: Validate License
	// TODO: Revoke License
	// TODO: register to basic auth
}
