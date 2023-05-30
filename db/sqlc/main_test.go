package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestQueries *Queries
var TestConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	TestConn, err = sql.Open("postgres", "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db %w", err)
	}

	TestQueries = New(TestConn)

	os.Exit(m.Run())
}
