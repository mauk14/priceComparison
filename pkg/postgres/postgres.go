package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenDb(dsn string) (*pgxpool.Pool, error) {

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return db, nil

}
