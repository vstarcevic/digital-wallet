package messaging

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"slices"

	"transaction-service/database"
	"transaction-service/model"

	"github.com/IBM/sarama"
)

func CreateTopicsIfNotExists(kafkaUrl string) {

	brokerAddress := []string{kafkaUrl}

	// create topics
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(brokerAddress, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	// Get list of topics
	topics, err := consumer.Topics()
	if err != nil {
		panic(err)
	}

	if !slices.Contains(topics, "user-created") {
		config.Version = sarama.V3_6_0_0
		admin, err := sarama.NewClusterAdmin(brokerAddress, config)
		if err != nil {
			log.Fatal("Error while creating cluster admin: ", err.Error())
		}
		defer func() {
			_ = admin.Close()
		}()
		err = admin.CreateTopic("user-created", &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			log.Fatal("Error while creating topic: ", err.Error())
		}
	}
}

func ListenTopic(topics string, conn *sql.DB, kafkaUrl string) {

	brokerAddress := []string{kafkaUrl}

	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokerAddress, config)
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topics, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var data model.User
			err := json.Unmarshal(msg.Value, &data)
			if err != nil {
				log.Printf("Consumer error, %s - value not supported.", msg.Value)
				continue
			}
			database.InsertBalance(conn, data)
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumer finished.")

}
