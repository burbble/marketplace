package db

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	maxRetries       = 10
	pingTimeout      = 3 * time.Second
	maxOpenConns     = 20
	maxIdleConns     = 20
	connMaxLifetime  = 30 * time.Minute
	connMaxIdleTime  = 5 * time.Minute
	initialBackoff   = 500 * time.Millisecond
	backoffMultiplier = 2
)

type PostgresConfig interface {
	GetHost() string
	GetPort() string
	GetUser() string
	GetPassword() string
	GetDBName() string
}

type Connection struct {
	DB      *sqlx.DB
	Builder sq.StatementBuilderType
}

func NewConnection(ctx context.Context, cfg PostgresConfig, lg *zap.Logger) (*Connection, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.GetHost(), cfg.GetPort(), cfg.GetUser(), cfg.GetPassword(), cfg.GetDBName(),
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	backoff := initialBackoff
	for i := 0; i < maxRetries; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
		err = db.PingContext(pingCtx)
		cancel()

		if err == nil {
			break
		}

		lg.Warn("postgres ping failed, retrying",
			zap.Int("attempt", i+1),
			zap.Duration("backoff", backoff),
			zap.Error(err),
		)

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while connecting to postgres: %w", ctx.Err())
		case <-time.After(backoff):
		}

		backoff *= backoffMultiplier
	}

	if err != nil {
		return nil, fmt.Errorf("postgres ping after %d retries: %w", maxRetries, err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if _, err := db.ExecContext(ctx, "SET timezone = 'UTC'"); err != nil {
		return nil, fmt.Errorf("set timezone: %w", err)
	}

	lg.Info("postgres connected",
		zap.String("host", cfg.GetHost()),
		zap.String("db", cfg.GetDBName()),
	)

	return &Connection{
		DB:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}

func (c *Connection) Close() {
	if c.DB != nil {
		_ = c.DB.Close()
	}
}
