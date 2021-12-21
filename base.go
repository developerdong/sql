package sql

import (
	"context"
	"database/sql"
)

// BaseDB is the most inner middleware, which implements the DB interface. Other
// middlewares can wrap each other casually, but this one is the base.
type BaseDB struct {
	*sql.DB
}

func (b *BaseDB) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmt, err := b.DB.PrepareContext(ctx, query)
	return &BaseStmt{stmt}, err
}

func (b *BaseDB) Prepare(query string) (Stmt, error) {
	stmt, err := b.DB.Prepare(query)
	return &BaseStmt{stmt}, err
}

func (b *BaseDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := b.DB.BeginTx(ctx, opts)
	return &BaseTx{tx}, err
}

func (b *BaseDB) Begin() (Tx, error) {
	tx, err := b.DB.Begin()
	return &BaseTx{tx}, err
}

func (b *BaseDB) Conn(ctx context.Context) (Conn, error) {
	conn, err := b.DB.Conn(ctx)
	return &BaseConn{conn}, err
}
func (b *BaseDB) OriginDB() *sql.DB {
	return b.DB
}

// BaseStmt implements the Stmt interface.
type BaseStmt struct {
	*sql.Stmt
}

func (b *BaseStmt) OriginStmt() *sql.Stmt {
	return b.Stmt
}

// BaseTx implements the Tx interface.
type BaseTx struct {
	*sql.Tx
}

func (b *BaseTx) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmt, err := b.Tx.PrepareContext(ctx, query)
	return &BaseStmt{stmt}, err
}
func (b *BaseTx) Prepare(query string) (Stmt, error) {
	stmt, err := b.Tx.Prepare(query)
	return &BaseStmt{stmt}, err
}
func (b *BaseTx) StmtContext(ctx context.Context, stmt Stmt) Stmt {
	return &BaseStmt{b.Tx.StmtContext(ctx, stmt.OriginStmt())}
}
func (b *BaseTx) Stmt(stmt Stmt) Stmt {
	return &BaseStmt{b.Tx.Stmt(stmt.OriginStmt())}
}

func (b *BaseTx) OriginTx() *sql.Tx {
	return b.Tx
}

// BaseConn implements the Conn interface.
type BaseConn struct {
	*sql.Conn
}

func (b *BaseConn) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmt, err := b.Conn.PrepareContext(ctx, query)
	return &BaseStmt{stmt}, err
}
func (b *BaseConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := b.Conn.BeginTx(ctx, opts)
	return &BaseTx{tx}, err
}
func (b *BaseConn) OriginConn() *sql.Conn {
	return b.Conn
}
