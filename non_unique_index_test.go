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

const (
	MEAN               = 2
	STANDARD_DEVIATION = 3
)

func BenchmarkRandomSelectsFromNonUniqueIndex(b *testing.B) {
	conn, err := sqlx.Connect("mysql", DSN)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()
	var count int
	err = conn.Get(&count, `SELECT COUNT(DISTINCT id) AS count FROM non_unique_index_table`)
	if err != nil {
		b.Fatal(err)
	}

	if count < SIZE {
		bootstrapNonUniqueIndexedTable(conn, SIZE-count)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := rand.Intn(SIZE) + 1
		row := row{}
		err := conn.Get(&row, `SELECT id, SUM(value) AS value, ANY_VALUE(metadata) AS metadata FROM non_unique_index_table WHERE id = ? GROUP BY id`, id)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func bootstrapNonUniqueIndexedTable(conn *sqlx.DB, more int) {
	rows := make([]row, 0, BATCH)
	conveyer := make(chan []row, WORKERS)
	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go func(conn *sqlx.DB, conveyer <-chan []row) {
			defer wg.Done()
			for rows := range conveyer {
				err := batch(conn, `INSERT INTO non_unique_index_table (id, metadata, value) VALUES (:id, :metadata, :value)`, rows)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(conn, conveyer)
	}

	var id int
	err := conn.Get(&id, `SELECT COALESCE(MAX(id), 0) as id FROM non_unique_index_table`)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < more; i++ {
		id++
		duplicates := int(math.Max(1, math.Round(rand.NormFloat64()*float64(STANDARD_DEVIATION)+float64(MEAN))))

		name := faker.FirstName() + " " + faker.LastName()

		for j := 0; j < duplicates; j++ {
			rows = append(rows, row{
				ID:       id,
				Value:    math.Round(rand.Float64()*1000) / 10,
				Metadata: name,
			})

			if len(rows) == BATCH {
				conveyer <- rows
				rows = make([]row, 0, BATCH)
			}
		}
	}

	if len(rows) > 0 {
		conveyer <- rows
	}

	close(conveyer)
	wg.Wait()
}
