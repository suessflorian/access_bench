package main

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
)

func BenchmarkRandomSelectsFromUniqueIndex(b *testing.B) {
	conn, err := sqlx.Connect("mysql", DSN)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()
	var count int
	err = conn.Get(&count, `SELECT COUNT(id) AS count FROM unique_index_table`)
	if err != nil {
		b.Fatal(err)
	}

	if count < SIZE {
		bootstrapUniqueIndexedTable(conn, SIZE-count)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := rand.Intn(SIZE) + 1
		row := row{}
		err := conn.Get(&row, `SELECT id, value, metadata FROM unique_index_table WHERE id = ?`, id)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func bootstrapUniqueIndexedTable(conn *sqlx.DB, more int) {
	rows := make([]row, 0, BATCH)
	conveyer := make(chan []row, WORKERS)
	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go func(conn *sqlx.DB, conveyer <-chan []row) {
			defer wg.Done()
			for rows := range conveyer {
				err := batch(conn, `INSERT INTO unique_index_table (metadata, value) VALUES (:metadata, :value)`, rows)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(conn, conveyer)
	}

	for i := 0; i < more; i++ {
		rows = append(rows, row{
			Metadata: faker.FirstName() + " " + faker.LastName(),
			Value:    math.Round(rand.Float64()*1000) / 10,
		})

		if len(rows) == BATCH {
			conveyer <- rows
			rows = make([]row, 0, BATCH)
		}
	}

	if len(rows) > 0 {
		conveyer <- rows
	}

	close(conveyer)
	wg.Wait()
}
