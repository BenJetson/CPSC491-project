package db

import "github.com/jmoiron/sqlx"

type DB struct {
	*sqlx.DB
	logger *
}
