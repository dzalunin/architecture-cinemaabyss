package broker

import (
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

func NewProducer(brokers string) sarama.SyncProducer {
	brokersSet := strings.Split(brokers, ",")
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 2
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokersSet, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	return producer
}
