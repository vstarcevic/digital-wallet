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
	postgresAddr := os.Getenv("POSTGRES_ADDR")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DB")
	postgresPort := os.Getenv("POSTGRES_PORT")

	natsAddr := os.Getenv("NATS_URL")
	natsPort := os.Getenv("NATS_PORT")

	kafkaAddr := os.Getenv("KAFKA_URL")
	kafkaPort := os.Getenv("KAFKA_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPass, postgresAddr, postgresPort, postgresDb)
	if postgresAddr == "" {
		dsn = "postgres://postgres:password@localhost:5432/userdb?sslmode=disable"
	}

	natsUrl := fmt.Sprintf("%s:%s", natsAddr, natsPort)
	if natsAddr == "" {
		natsUrl = "localhost:4222"
	}

	kafkaUrl := fmt.Sprintf("%s:%s", kafkaAddr, kafkaPort)
	if kafkaAddr == "" {
		kafkaUrl = "localhost:9092"
	}

	// connect and open db
	dbConn := database.ConnectToDB(dsn)

	// run all migrations
	database.RunMigrations(dsn)
	defer dbConn.Close()

	// create topic in Kafka
	// this would not be in production,
	// we would be creating manually all topics needed.
	messaging.CreateTopicsIfNotExists(kafkaUrl)

	// connect to nats
	natsConn := messaging.ConnectToNats(natsUrl)
	defer natsConn.Close()

	cfg := api.Config{
		Db:  dbConn,
		Nts: natsConn,
		App: api.AppSettings{
			Dsn:      dsn,
			KafkaUrl: kafkaUrl,
			NatsUrl:  natsUrl,
		},
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
