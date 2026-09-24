package db

import "database/sql"

type Store struct {
	*Queries
}

func NewStore(conn *sql.DB) *Store {
	return &Store{
		Queries: New(conn),
	}
}