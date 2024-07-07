package messaging

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/IBM/sarama"
)

type SaramaProducer struct {
	producer *sarama.SyncProducer
}

func CreateTopicsIfNotExists() {

	// create topics
	brokerAddrs := []string{"localhost:9092"}
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(brokerAddrs, config)
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
		admin, err := sarama.NewClusterAdmin(brokerAddrs, config)
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

func PublishJSON(topic string, data any, host string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonData),
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{host}, config)
	if err != nil {
		return err
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
