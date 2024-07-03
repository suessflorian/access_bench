package main

import "github.com/jmoiron/sqlx"

const (
	BATCH   = 10_000
	WORKERS = 100
	SIZE    = 10_000_000
)

type row struct {
	ID       int
	Value    float64
	Metadata string
}

func batch[T any](conn *sqlx.DB, prepared string, load []T) error {
	tx, err := conn.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(prepared)
	if err != nil {
		return err
	}

	for _, row := range load {
		_, err := stmt.Exec(row)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
