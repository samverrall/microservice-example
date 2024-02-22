package postgresutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MigrateFS(ctx context.Context, db *pgxpool.Pool, fsys fs.FS) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("rollback tx: %w", err)
		}
	}()

	migrationsDir, err := fs.Sub(fsys, "migrations")
	if err != nil {
		return fmt.Errorf("sub migrations dir: %w", err)
	}

	files, err := fs.ReadDir(migrationsDir, ".")
	if err != nil {
		return fmt.Errorf("read migration dir: %w", err)
	}

	queries := make([]string, len(files))
	for i, f := range files {
		filename := f.Name()

		if f.IsDir() {
			return fmt.Errorf("want file; got directory: %q", filename)
		}

		if filename[:4] != fmt.Sprintf("%04d", i+1) {
			return fmt.Errorf("want filename to start with %04d; got %q", i+1, filename)
		}

		b, err := fs.ReadFile(migrationsDir, filename)
		if err != nil {
			return fmt.Errorf("read migration file: %w", err)
		}

		queries[i] = string(b)
	}

	if err := migrate(ctx, tx, queries); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func migrate(ctx context.Context, tx pgx.Tx, migrations []string) error {
	if len(migrations) == 0 {
		return nil
	}

	for _, stmt := range migrations {
		if stmt = strings.TrimSpace(stmt); stmt == "" {
			continue
		}
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("exec migration: %w", err)
		}
	}

	return nil
}
