package postgres

import (
	"context"
	"fmt"

	"github.com/mfajri11/family-catering-micro-service/auth-service/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultAttemps = 5

func MustNew() *pgxpool.Pool {
	var err error
	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.Cfg.Postgres.ConnString())
	if err != nil {
		panic(err)
	}
	for i := 1; i <= defaultAttemps; i++ {
		err = db.Ping(ctx)
		if err == nil {
			continue
		}
		fmt.Printf("postgres.MustNew: fail ping to database, left attempts  %d", i+1)
	}

	if err != nil {
		panic(err)
	}

	return db
}
