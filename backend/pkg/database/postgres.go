package database

import (
    "context"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/rs/zerolog"
)

type DB struct {
    *sqlx.DB
    logger *zerolog.Logger
}

type Config struct {
    Host         string
    Port         string
    User         string
    Password     string
    Database     string
    SSLMode      string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}

func NewPostgresDB(cfg Config, logger *zerolog.Logger) (*DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
    )

    db, err := sqlx.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Configure connection pool
    if cfg.MaxOpenConns == 0 {
        cfg.MaxOpenConns = 25
    }
    if cfg.MaxIdleConns == 0 {
        cfg.MaxIdleConns = 5
    }
    if cfg.MaxLifetime == 0 {
        cfg.MaxLifetime = 5 * time.Minute
    }

    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.MaxLifetime)

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    logger.Info().
        Str("host", cfg.Host).
        Str("database", cfg.Database).
        Msg("Connected to PostgreSQL database")

    return &DB{
        DB:     db,
        logger: logger,
    }, nil
}

func (db *DB) Close() error {
    if err := db.DB.Close(); err != nil {
        db.logger.Error().Err(err).Msg("Failed to close database connection")
        return err
    }
    db.logger.Info().Msg("Database connection closed")
    return nil
}

// Transaction wrapper with automatic rollback on error
func (db *DB) Transaction(fn func(*sqlx.Tx) error) error {
    tx, err := db.Beginx()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
    }()

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("tx error: %w, rollback error: %v", err, rbErr)
        }
        return err
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

// Migrate runs database migrations
func (db *DB) Migrate(migrationsPath string) error {
    // This would use golang-migrate in production
    // For now, we'll handle migrations through Docker or manually
    db.logger.Info().Msg("Database migrations should be run separately")
    return nil
}

// Helper methods
func (db *DB) Exists(query string, args ...interface{}) (bool, error) {
    var exists bool
    err := db.Get(&exists, query, args...)
    return exists, err
}

func (db *DB) Count(query string, args ...interface{}) (int, error) {
    var count int
    err := db.Get(&count, query, args...)
    return count, err
}

