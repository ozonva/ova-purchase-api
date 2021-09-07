package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type DB struct {
	Db *sqlx.DB
}

func NewDB(url string) (*DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	newDb := sqlx.NewDb(db, "pgx")
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}
	return &DB{
		Db: newDb,
	}, nil
}

func (s *DB) Disposal() {
	log.Debug().Msg("Closing db...")
	if err := s.Db.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close db")
	}
}
