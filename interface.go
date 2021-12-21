package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

type DB interface {
	PingContext(ctx context.Context) error
	Ping() error
	Close() error
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	Stats() sql.DBStats
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	Prepare(query string) (Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Begin() (Tx, error)
	Driver() driver.Driver
	Conn(ctx context.Context) (Conn, error)
	OriginDB() *sql.DB
}

type Stmt interface {
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	Exec(args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
	Query(args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	QueryRow(args ...interface{}) *sql.Row
	Close() error
	OriginStmt() *sql.Stmt
}

type Tx interface {
	Commit() error
	Rollback() error
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	Prepare(query string) (Stmt, error)
	StmtContext(ctx context.Context, stmt Stmt) Stmt
	Stmt(stmt Stmt) Stmt
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	OriginTx() *sql.Tx
}

type Conn interface {
	PingContext(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	Raw(f func(driverConn interface{}) error) (err error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Close() error
	OriginConn() *sql.Conn
}
