package store

import "database/sql"

type Driver interface {
	Close() error
	GetDB() *sql.DB
	Migrate() error
}
