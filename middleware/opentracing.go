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

type Tracer struct {
	sql.DB
}

func (t *Tracer) PingContext(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PingContext")
	defer span.Finish()
	err := t.DB.PingContext(ctx)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}

func (t *Tracer) Ping() error {
	span := opentracing.StartSpan("Ping")
	defer span.Finish()
	err := t.DB.Ping()
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}

func (t *Tracer) Close() error {
	span := opentracing.StartSpan("Close")
	defer span.Finish()
	err := t.DB.Close()
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return err
}

func (t *Tracer) SetMaxIdleConns(n int) {
	span := opentracing.StartSpan("SetMaxIdleConns")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Int("n", n))
	t.DB.SetMaxIdleConns(n)
}

func (t *Tracer) SetMaxOpenConns(n int) {
	span := opentracing.StartSpan("SetMaxOpenConns")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Int("n", n))
	t.DB.SetMaxOpenConns(n)
}

func (t *Tracer) SetConnMaxLifetime(d time.Duration) {
	span := opentracing.StartSpan("SetConnMaxLifetime")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("d", d.String()))
	t.DB.SetConnMaxLifetime(d)
}

func (t *Tracer) SetConnMaxIdleTime(d time.Duration) {
	span := opentracing.StartSpan("SetConnMaxIdleTime")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("d", d.String()))
	t.DB.SetConnMaxIdleTime(d)
}

func (t *Tracer) Stats() stdSql.DBStats {
	span := opentracing.StartSpan("Stats")
	defer span.Finish()
	return t.DB.Stats()
}

func (t *Tracer) PrepareContext(ctx context.Context, query string) (*stdSql.Stmt, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PrepareContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.DB.PrepareContext(ctx, query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return stmt, err
}

func (t *Tracer) Prepare(query string) (*stdSql.Stmt, error) {
	span := opentracing.StartSpan("Prepare")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query))
	stmt, err := t.DB.Prepare(query)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return stmt, err
}

func (t *Tracer) ExecContext(ctx context.Context, query string, args ...interface{}) (stdSql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExecContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}

func (t *Tracer) Exec(query string, args ...interface{}) (stdSql.Result, error) {
	span := opentracing.StartSpan("Exec")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	result, err := t.DB.Exec(query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return result, err
}

func (t *Tracer) QueryContext(ctx context.Context, query string, args ...interface{}) (*stdSql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}

func (t *Tracer) Query(query string, args ...interface{}) (*stdSql.Rows, error) {
	span := opentracing.StartSpan("Query")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	rows, err := t.DB.Query(query, args...)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return rows, err
}

func (t *Tracer) QueryRowContext(ctx context.Context, query string, args ...interface{}) *stdSql.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryRowContext")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.DB.QueryRowContext(ctx, query, args...)
}

func (t *Tracer) QueryRow(query string, args ...interface{}) *stdSql.Row {
	span := opentracing.StartSpan("QueryRow")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.String("query", query), log.Object("args", args))
	return t.DB.QueryRow(query, args...)
}

func (t *Tracer) BeginTx(ctx context.Context, opts *stdSql.TxOptions) (*stdSql.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BeginTx")
	defer span.Finish()
	span.LogFields(log.Event("debug"), log.Object("opts", opts))
	tx, err := t.DB.BeginTx(ctx, opts)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return tx, err
}

func (t *Tracer) Begin() (*stdSql.Tx, error) {
	span := opentracing.StartSpan("Begin")
	defer span.Finish()
	tx, err := t.DB.Begin()
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return tx, err
}

func (t *Tracer) Driver() driver.Driver {
	span := opentracing.StartSpan("Driver")
	defer span.Finish()
	return t.DB.Driver()
}

func (t *Tracer) Conn(ctx context.Context) (*stdSql.Conn, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Conn")
	defer span.Finish()
	conn, err := t.DB.Conn(ctx)
	if err != nil {
		span.LogFields(log.Event("error"), log.Error(err))
	}
	return conn, err
}
