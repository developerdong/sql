package middleware

import (
	"context"
	stdSql "database/sql"
	"database/sql/driver"
	"github.com/developerdong/sql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

var (
	_ sql.DB   = (*TraceDB)(nil)
	_ sql.Stmt = (*TraceStmt)(nil)
	_ sql.Tx   = (*TraceTx)(nil)
	_ sql.Conn = (*TraceConn)(nil)
)

type TraceDB struct {
	sql.DB
}

func (t *TraceDB) PingContext(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PingContext")
	defer span.Finish()
	err := t.DB.PingContext(ctx)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}

func (t *TraceDB) Ping() error {
	span := opentracing.StartSpan("Ping")
	defer span.Finish()
	err := t.DB.Ping()
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}

func (t *TraceDB) Close() error {
	span := opentracing.StartSpan("Close")
	defer span.Finish()
	err := t.DB.Close()
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}

func (t *TraceDB) SetMaxIdleConns(n int) {
	span := opentracing.StartSpan("SetMaxIdleConns")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Int("n", n))
	t.DB.SetMaxIdleConns(n)
}

func (t *TraceDB) SetMaxOpenConns(n int) {
	span := opentracing.StartSpan("SetMaxOpenConns")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Int("n", n))
	t.DB.SetMaxOpenConns(n)
}

func (t *TraceDB) SetConnMaxLifetime(d time.Duration) {
	span := opentracing.StartSpan("SetConnMaxLifetime")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("d", d.String()))
	t.DB.SetConnMaxLifetime(d)
}

func (t *TraceDB) SetConnMaxIdleTime(d time.Duration) {
	span := opentracing.StartSpan("SetConnMaxIdleTime")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("d", d.String()))
	t.DB.SetConnMaxIdleTime(d)
}

func (t *TraceDB) Stats() stdSql.DBStats {
	span := opentracing.StartSpan("Stats")
	defer span.Finish()
	return t.DB.Stats()
}

func (t *TraceDB) PrepareContext(ctx context.Context, query string) (sql.Stmt, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PrepareContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.DB.PrepareContext(ctx, query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceStmt{stmt, query}, err
}

func (t *TraceDB) Prepare(query string) (sql.Stmt, error) {
	span := opentracing.StartSpan("Prepare")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.DB.Prepare(query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceStmt{stmt, query}, err
}

func (t *TraceDB) ExecContext(ctx context.Context, query string, args ...interface{}) (stdSql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExecContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}

func (t *TraceDB) Exec(query string, args ...interface{}) (stdSql.Result, error) {
	span := opentracing.StartSpan("Exec")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.DB.Exec(query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}

func (t *TraceDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*stdSql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}

func (t *TraceDB) Query(query string, args ...interface{}) (*stdSql.Rows, error) {
	span := opentracing.StartSpan("Query")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.DB.Query(query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}

func (t *TraceDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *stdSql.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryRowContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.DB.QueryRowContext(ctx, query, args...)
}

func (t *TraceDB) QueryRow(query string, args ...interface{}) *stdSql.Row {
	span := opentracing.StartSpan("QueryRow")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.DB.QueryRow(query, args...)
}

func (t *TraceDB) BeginTx(ctx context.Context, opts *stdSql.TxOptions) (sql.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BeginTx")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Object("opts", opts))
	tx, err := t.DB.BeginTx(ctx, opts)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceTx{tx}, err
}

func (t *TraceDB) Begin() (sql.Tx, error) {
	span := opentracing.StartSpan("Begin")
	defer span.Finish()
	tx, err := t.DB.Begin()
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceTx{tx}, err
}

func (t *TraceDB) Driver() driver.Driver {
	span := opentracing.StartSpan("Driver")
	defer span.Finish()
	return t.DB.Driver()
}

func (t *TraceDB) Conn(ctx context.Context) (sql.Conn, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Conn")
	defer span.Finish()
	conn, err := t.DB.Conn(ctx)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceConn{conn}, err
}

type TraceStmt struct {
	sql.Stmt
	query string
}

func (s *TraceStmt) ExecContext(ctx context.Context, args ...interface{}) (stdSql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExecContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", s.query), log.Object("args", args))
	result, err := s.Stmt.ExecContext(ctx, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}
func (s *TraceStmt) Exec(args ...interface{}) (stdSql.Result, error) {
	span := opentracing.StartSpan("Exec")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", s.query), log.Object("args", args))
	result, err := s.Stmt.Exec(args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}
func (s *TraceStmt) QueryContext(ctx context.Context, args ...interface{}) (*stdSql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", s.query), log.Object("args", args))
	result, err := s.Stmt.QueryContext(ctx, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}
func (s *TraceStmt) Query(args ...interface{}) (*stdSql.Rows, error) {
	span := opentracing.StartSpan("Query")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", s.query), log.Object("args", args))
	result, err := s.Stmt.Query(args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}
func (s *TraceStmt) QueryRowContext(ctx context.Context, args ...interface{}) *stdSql.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryRowContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", s.query), log.Object("args", args))
	return s.Stmt.QueryRowContext(ctx, args...)
}
func (s *TraceStmt) QueryRow(args ...interface{}) *stdSql.Row {
	span := opentracing.StartSpan("QueryRow")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", s.query), log.Object("args", args))
	return s.Stmt.QueryRow(args...)
}

type TraceTx struct {
	sql.Tx
}

func (t *TraceTx) PrepareContext(ctx context.Context, query string) (sql.Stmt, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PrepareContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.Tx.PrepareContext(ctx, query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceStmt{stmt, query}, err
}

func (t *TraceTx) Prepare(query string) (sql.Stmt, error) {
	span := opentracing.StartSpan("Prepare")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.Tx.Prepare(query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceStmt{stmt, query}, err
}

func (t *TraceTx) StmtContext(ctx context.Context, stmt sql.Stmt) sql.Stmt {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StmtContext")
	defer span.Finish()
	return &TraceStmt{t.Tx.StmtContext(ctx, stmt), ""}
}

func (t *TraceTx) Stmt(stmt sql.Stmt) sql.Stmt {
	span := opentracing.StartSpan("Stmt")
	defer span.Finish()
	return &TraceStmt{t.Tx.Stmt(stmt), ""}
}

func (t *TraceTx) ExecContext(ctx context.Context, query string, args ...interface{}) (stdSql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExecContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.Tx.ExecContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}

func (t *TraceTx) Exec(query string, args ...interface{}) (stdSql.Result, error) {
	span := opentracing.StartSpan("Exec")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.Tx.Exec(query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}

func (t *TraceTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*stdSql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}

func (t *TraceTx) Query(query string, args ...interface{}) (*stdSql.Rows, error) {
	span := opentracing.StartSpan("Query")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.Tx.Query(query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}

func (t *TraceTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *stdSql.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryRowContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.Tx.QueryRowContext(ctx, query, args...)
}

func (t *TraceTx) QueryRow(query string, args ...interface{}) *stdSql.Row {
	span := opentracing.StartSpan("QueryRow")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.Tx.QueryRow(query, args...)
}

type TraceConn struct {
	sql.Conn
}

func (t *TraceConn) PingContext(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PingContext")
	defer span.Finish()
	err := t.Conn.PingContext(ctx)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}
func (t *TraceConn) ExecContext(ctx context.Context, query string, args ...interface{}) (stdSql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExecContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}
func (t *TraceConn) QueryContext(ctx context.Context, query string, args ...interface{}) (*stdSql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}
func (t *TraceConn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *stdSql.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryRowContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.Conn.QueryRowContext(ctx, query, args...)
}
func (t *TraceConn) PrepareContext(ctx context.Context, query string) (sql.Stmt, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PrepareContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceStmt{stmt, query}, err
}
func (t *TraceConn) BeginTx(ctx context.Context, opts *stdSql.TxOptions) (sql.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BeginTx")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Object("opts", opts))
	tx, err := t.Conn.BeginTx(ctx, opts)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return &TraceTx{tx}, err
}
