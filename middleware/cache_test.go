package middleware

import (
	"context"
	"database/sql"
	s "github.com/developerdong/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
	"testing"
)

const (
	// Dsn is the data source name.
	Dsn = "root:root@tcp(localhost:3306)/mysql"
	// N is the loop times.
	N = 100000
	// MaxOpenConns is the concurrent connections limit.
	MaxOpenConns = 10
	// Query is the query statement for test.
	Query = "SELECT TRUE FROM user WHERE user=?;"
	// Arg is the query argument for test.
	Arg = "root"
)

func TestCacheDB_QueryRowContext_Serial(t *testing.T) {
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		t.Fatal("can not open the connection to mysql", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	db.SetMaxOpenConns(1)
	cacheDb := CacheDB{DB: &s.BaseDB{DB: db}}
	var result int
	for i := 0; i < N; i++ {
		if err := cacheDb.QueryRowContext(context.Background(), Query, Arg).Scan(&result); err != nil {
			t.Error(err)
		}
	}
}

func TestDB_QueryRowContext_Serial(t *testing.T) {
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		t.Fatal("can not open the connection to mysql", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	db.SetMaxOpenConns(1)
	var result int
	for i := 0; i < N; i++ {
		if err := db.QueryRowContext(context.Background(), Query, Arg).Scan(&result); err != nil {
			t.Error(err)
		}
	}
}

func TestCacheDB_QueryRowContext_Concurrent(t *testing.T) {
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		t.Fatal("can not open the connection to mysql", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	db.SetMaxOpenConns(MaxOpenConns)
	cacheDb := CacheDB{DB: &s.BaseDB{DB: db}}
	g, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < N; i++ {
		g.Go(func() error {
			var result int
			return cacheDb.QueryRowContext(ctx, Query, Arg).Scan(&result)
		})
	}
	if err := g.Wait(); err != nil {
		t.Error(err)
	}
}

func TestDB_QueryRowContext_Concurrent(t *testing.T) {
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		t.Fatal("can not open the connection to mysql", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	db.SetMaxOpenConns(MaxOpenConns)
	g, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < N; i++ {
		g.Go(func() error {
			var result int
			return db.QueryRowContext(ctx, Query, Arg).Scan(&result)
		})
	}
	if err := g.Wait(); err != nil {
		t.Error(err)
	}
}
