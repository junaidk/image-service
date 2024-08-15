package postgress

import (
	"context"
	"database/sql"
	"fmt"

	imageapi "github.com/junaidk/image-service"
	_ "github.com/lib/pq"
)

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()

	DSN string
}

func NewDB(dsn string) *DB {

	db := &DB{
		DSN: dsn,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	if db.db, err = sql.Open("postgres", db.DSN); err != nil {
		return err
	}

	return db.db.Ping()
}

func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

type Tx struct {
	*sql.Tx
	db *DB
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx: tx,
		db: db,
	}, nil
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func FormatError(err error) error {
	if err == nil {
		return nil
	}

	switch err.Error() {
	case "sql: no rows in result set":
		return imageapi.Errorf(imageapi.ERRNOTFOUND, "Resource not found.")
	default:
		return err
	}
}
