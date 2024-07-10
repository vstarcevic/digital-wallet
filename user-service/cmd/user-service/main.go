package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"user-service/api"
	"user-service/database"
	"user-service/messaging"
)

func main() {
	dsn := os.Getenv("DSN_USER_DB")

	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5432/userdb?sslmode=disable"
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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "9000"),
		Handler: api.Routes(&cfg),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
