package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/MikoBerries/SimpleBank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	//load config file using viper
	cf, err := util.LoadConfig("../..")
	if err != nil {
		log.Panic(err)
	}

	testConn, err = sql.Open(cf.DBDriver, cf.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db %w", err)
	}

	testQueries = New(testConn)

	os.Exit(m.Run())
}
