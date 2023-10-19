package ports

import (
	"context"

	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/domain"
)

type IAuthRepository interface {
	// InsertSession insert a session to a database and cache it to redis
	// SessionByEmail()
	// SessionByXX()
	InsertSession(ctx context.Context) error
	Session(ctx context.Context, sid string) (session domain.Session, err error)
	SessionIDByEmail(ctx context.Context, email string) (string, error)
}

type IAuthCache interface {
	// SetSession
	// (Get)Session
	// Revoke

}
