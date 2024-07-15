package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store interface {
	Querier
}

type SqlStore struct {
	conn *sql.DB
	*Queries
}

func NewStore() (Store, error) {
	db, err := sql.Open("postgres", "postgresql://root:postgres@localhost:5432/webSocketDB?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &SqlStore{
		conn:    db,
		Queries: New(db),
	}, nil
}
