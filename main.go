package main

import (
	"assignment_01/simplebank/api"
	"assignment_01/simplebank/util"

	"database/sql"

	db "assignment_01/simplebank/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load configration:", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("can not connect to db %v ", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err = server.Start(config.Addres); err != nil {
		log.Fatal("can not start server :", err)
	}

}
