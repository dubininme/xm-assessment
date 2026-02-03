package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dubininme/xm-assessment/internal/config"
	"github.com/dubininme/xm-assessment/internal/delivery/http/handler"
	"github.com/dubininme/xm-assessment/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Db struct {
	*sql.DB
}

func Connect(ctx context.Context, cfg config.DbConfig) (*Db, error) {
	log := logger.FromContext(ctx)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Second)

	log.Info("connected to database successfully",
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"database", cfg.DBName)

	return &Db{db}, nil
}

var _ handler.HealthChecker = (*DBHealthChecker)(nil)

type DBHealthChecker struct {
	db *Db
}

func NewDBHealthChecker(db *Db) *DBHealthChecker {
	return &DBHealthChecker{db: db}
}

func (c *DBHealthChecker) Check(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *DBHealthChecker) Name() string {
	return "postgres"
}
