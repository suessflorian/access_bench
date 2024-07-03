package main

import (
	"context"
	"log"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const DSN = "user:password@tcp(127.0.0.1:3306)/main?multiStatements=true"

//go:embed test.sql
var migration string

func init() {
	ctx := context.Background()

	conn, err := sqlx.Open("mysql", DSN)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	_, err = conn.ExecContext(ctx, migration)
	if err != nil {
		log.Fatalf("error migrating up: %v", err)
	}
}
