package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/safe-homie/backend/internal/config"
)

type _postgres struct {
	db  *pgxpool.Pool
	cfg *config.StoreConfig
}

func New(cfg *config.Config) (*_postgres, error) {
	if cfg.Store.DSN == "" {
		return nil, errors.New("dsn required")
	}
	db, err := pgxpool.New(context.Background(), cfg.Store.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open db with dsn: %w", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &_postgres{db: db, cfg: &cfg.Store}, nil
}

func (p *_postgres) GetDB() *pgxpool.Pool {
	return p.db
}

func (p *_postgres) Close() {
	p.db.Close()
}

func (p *_postgres) Migrate() error {
	sqlDB, err := sql.Open(p.cfg.Driver, p.cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to open db for migration: %w", err)
	}
	defer sqlDB.Close()
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
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
	// if err := m.Force(0); err != nil {
	// 	return fmt.Errorf("failed to reset migration: %w", err)
	// }
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
