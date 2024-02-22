package postgresutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildDSN(host, user, password, dbname string, port int) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", host, user, password, dbname, port)
}

func Connect(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return db, nil
}
