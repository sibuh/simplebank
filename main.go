package main

import (
	"exercise/simplebank/api"
	"exercise/simplebank/util"

	"database/sql"

	db "exercise/simplebank/db/sqlc"
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("failed to create server", err)
	}
	if err = server.Start(config.Addres); err != nil {
		log.Fatal("can not start server :", err)
	}

}
