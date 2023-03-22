package db

import (
	"database/sql"
	"exercise/simplebank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("can not load configration:", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("can not connect to db %v ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
