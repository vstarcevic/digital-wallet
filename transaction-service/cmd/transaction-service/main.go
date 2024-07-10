package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"transaction-service/api"
	"transaction-service/database"
	"transaction-service/messaging"
)

func main() {
	dsn := os.Getenv("DSN_BALANCE_DB")

	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5433/balancedb?sslmode=disable"
	}

	// connect and open db
	dbConn := database.ConnectToDB(dsn)

	// run all migrations
	database.RunMigrations(dsn)
	defer dbConn.Close()

	// create topic in Kafka
	// this would not be in production,
	// we would be creating manually all topics needed.
	messaging.CreateTopicsIfNotExists()

	cfg := api.Config{
		Db: dbConn,
	}

	go messaging.ListenTopic("user-created", cfg.Db)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "9001"),
		Handler: api.Routes(&cfg),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
