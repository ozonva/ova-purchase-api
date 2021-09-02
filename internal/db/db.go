package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(url string) (*sqlx.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	newDb := sqlx.NewDb(db, "pgx")
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}
	return newDb, nil
}
