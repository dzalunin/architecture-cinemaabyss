package broker

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
}

func NewConsumer(brokers string) *KafkaConsumer {
	brokersList := strings.Split(brokers, ",")
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokersList, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	return &KafkaConsumer{consumer: consumer}
}

func (kc *KafkaConsumer) Consume(ctx context.Context, topic string, partition int32, handler func(*sarama.ConsumerMessage)) {
	partitionConsumer, err := kc.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	log.Printf("Consuming from topic %s, partition %d", topic, partition)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return
		case msg := <-partitionConsumer.Messages():
			handler(msg)
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}

func (kc *KafkaConsumer) Close() {
	if err := kc.consumer.Close(); err != nil {
		log.Printf("Failed to close consumer: %v", err)
	}
}
