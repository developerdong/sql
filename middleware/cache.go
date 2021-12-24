package middleware

import (
	"context"
	stdSql "database/sql"
	"github.com/developerdong/sql"
	"strings"
	"sync"
)

var (
	_ sql.DB = (*CacheDB)(nil)
)

type CacheDB struct {
	sql.DB
	mu    sync.RWMutex
	stmts map[string]sql.Stmt
}

func (c *CacheDB) normalize(query string) string {
	query = strings.TrimSpace(query)
	if !strings.HasSuffix(query, ";") {
		query += ";"
	}
	return query
}

func (c *CacheDB) getStmt(query string) (stmt sql.Stmt, err error) {
	if c.stmts == nil {
		c.mu.Lock()
		if c.stmts == nil {
			c.stmts = make(map[string]sql.Stmt)
		}
		c.mu.Unlock()
	}
	query = c.normalize(query)
	c.mu.RLock()
	stmt = c.stmts[query]
	c.mu.RUnlock()
	if stmt == nil {
		c.mu.Lock()
		stmt = c.stmts[query]
		if stmt == nil {
			stmt, err = c.Prepare(query)
			if err == nil {
				c.stmts[query] = stmt
			}
		}
		c.mu.Unlock()
	}
	return
}

func (c *CacheDB) rmStmt(query string) {
	if c.stmts == nil {
		return
	}
	query = c.normalize(query)
	c.mu.Lock()
	stmt := c.stmts[query]
	delete(c.stmts, query)
	c.mu.Unlock()
	if stmt != nil {
		_ = stmt.Close()
	}
}

func (c *CacheDB) ExecContext(ctx context.Context, query string, args ...interface{}) (stdSql.Result, error) {
	stmt, err := c.getStmt(query)
	if err != nil {
		return c.DB.ExecContext(ctx, query, args...)
	}
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		c.rmStmt(query)
	}
	return result, err
}

func (c *CacheDB) Exec(query string, args ...interface{}) (stdSql.Result, error) {
	stmt, err := c.getStmt(query)
	if err != nil {
		return c.DB.Exec(query, args...)
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		c.rmStmt(query)
	}
	return result, err
}
func (c *CacheDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*stdSql.Rows, error) {
	stmt, err := c.getStmt(query)
	if err != nil {
		return c.DB.QueryContext(ctx, query, args...)
	}
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		c.rmStmt(query)
	}
	return rows, err
}
func (c *CacheDB) Query(query string, args ...interface{}) (*stdSql.Rows, error) {
	stmt, err := c.getStmt(query)
	if err != nil {
		return c.DB.Query(query, args...)
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		c.rmStmt(query)
	}
	return rows, err
}
func (c *CacheDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *stdSql.Row {
	stmt, err := c.getStmt(query)
	if err != nil {
		return c.DB.QueryRowContext(ctx, query, args...)
	}
	row := stmt.QueryRowContext(ctx, args...)
	if row.Err() != nil {
		c.rmStmt(query)
	}
	return row
}
func (c *CacheDB) QueryRow(query string, args ...interface{}) *stdSql.Row {
	stmt, err := c.getStmt(query)
	if err != nil {
		return c.DB.QueryRow(query, args...)
	}
	row := stmt.QueryRow(args...)
	if row.Err() != nil {
		c.rmStmt(query)
	}
	return row
}
