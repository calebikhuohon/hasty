package storage

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlDatetimeFormat = "2006-01-02 15:04:05.000000"
)

// scanner provides an interface to scan fields on *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...interface{}) error
}

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Ping() error
}

type Storage struct {
	db DB
}

func (s *Storage) PingDB() error {
	return s.db.Ping()
}

func New(db DB) *Storage {
	return &Storage{
		db: db,
	}
}

func isDuplicateErr(err error) bool {
	return err != nil && err.(*mysql.MySQLError).Number == 1062
}
