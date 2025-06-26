package broker

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

func NewProducerMessage(topic string, key string, payload any) (*sarama.ProducerMessage, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(bytes),
	}, nil
}
