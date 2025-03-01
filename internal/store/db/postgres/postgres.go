package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/safe-homie/backend/internal/config"
)

type _postgres struct {
	db  *sql.DB
	cfg *config.StoreConfig
}

func New(cfg *config.Config) (*_postgres, error) {
	if cfg.Store.DSN == "" {
		return nil, errors.New("dsn required")
	}
	db, err := sql.Open(cfg.Store.Driver, cfg.Store.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open db with dsn: %w", err)
	}
	return &_postgres{db: db, cfg: &cfg.Store}, nil
}

func (p *_postgres) GetDB() *sql.DB {
	return p.db
}

func (p *_postgres) Close() error {
	return p.db.Close()
}

func (p *_postgres) Migrate() error {
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("fail to create driver from db instance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migration/postgres",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("fail to create migration instance: %w", err)
	}
	if err := m.Force(0); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
