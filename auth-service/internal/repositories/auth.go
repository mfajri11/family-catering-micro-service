package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/domain"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/log"
	"github.com/redis/go-redis/v9"
)

const (
	sidUserKey string = "sid:user:%s"
)

type AuthRepository struct {
	db    *pgxpool.Pool
	cache *redis.Client
}

func New(db *pgxpool.Pool, cache *redis.Client) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) InsertSession(ctx context.Context, rows domain.Session) (sid string, err error) {
	// insert to db
	err = repo.db.QueryRow(ctx, insertSessionData).Scan(&sid)
	if err != nil {
		return "", err // TODO: wrap error
	}

	_, err = repo.cache.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.SetEx(ctx, fmt.Sprintf("email:sid:%s", rows.Email), sid, 1*time.Minute) // TODO: replace time.Minute
		p.HSet(ctx, fmt.Sprintf("sid:user:%s", sid), &rows)
		p.Expire(ctx, fmt.Sprintf("sid:user:%s", sid), 1*time.Minute) // TODO: replace time.Minute
		_, err := p.Exec(ctx)
		return err
	})

	if err != nil {
		// redis errors are ignored for now
		log.Error(err, "error cache the session") // TODO: wrap error
	}

	return sid, nil
}

func (repo *AuthRepository) SessionID(ctx context.Context, sid string) (*domain.Session, error) {
	key := fmt.Sprintf(sidUserKey, sid)
	val := &domain.Session{}
	err := repo.cache.Get(ctx, key).Scan(&val)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	if val == nil {
		// TODO: get from persistence database
		err := repo.db.QueryRow(ctx, selectSessionData, sid).Scan(&val)
		if err != nil {
			return nil, err // TODO: wrap error
		}
		return nil, nil
	}

	return val, nil
}

func (repo *AuthRepository) SessionIDByEmail(ctx context.Context, email string) (string, error) {
	var sid string
	err := repo.cache.Get(ctx, fmt.Sprintf("email:sid:%s", email)).Scan(&sid)
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return sid, nil
}

// func (repo *AuthRepository) redisSessionID(ctx context.Context, sid string, val any) error {
// 	key := fmt.Sprintf(sidUserKey, sid)
// 	err := repo.cache.HGetAll(ctx, key).Scan(&val)
// 	if errors.Is(err, redis.Nil) {
// 		return nil
// 	}
// 	if err != nil {
// 		return err // TODO: wrap error
// 	}

// 	return err
// }
