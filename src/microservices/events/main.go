package main

import (
	"context"
	"events/broker"
	"events/config"
	"events/dto"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	producer := broker.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	consumer := broker.NewConsumer(cfg.KafkaBrokers)
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go consumer.Consume(ctx, cfg.MovieEventsTopic, 0, messageHandler)
	go consumer.Consume(ctx, cfg.PaymentEventsTopic, 0, messageHandler)
	go consumer.Consume(ctx, cfg.UserEventsTopic, 0, messageHandler)

	router := gin.Default()

	eventsApi := router.Group("/api/events")
	{
		eventsApi.GET("/health", health)
		eventsApi.POST("/movie", movieEvent(producer, cfg.MovieEventsTopic))
		eventsApi.POST("/user", userEvent(producer, cfg.UserEventsTopic))
		eventsApi.POST("/payment", paymentEvent(producer, cfg.PaymentEventsTopic))
	}
	router.Run(":" + cfg.Port)
}

func health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": true})
}

func movieEvent(producer sarama.SyncProducer, topic string) gin.HandlerFunc {
	return func(context *gin.Context) {

		var movie dto.MovieEvent
		if err := context.ShouldBindJSON(&movie); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		event := dto.Event{
			ID:        fmt.Sprintf("movie-%d-%s", movie.MovieID, movie.Action),
			Type:      "movie",
			Timestamp: time.Now().UTC(),
			Payload:   movie,
		}

		msg, err := broker.NewProducerMessage(topic, event.Type, event)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := dto.EventResponse{
			Status:    "success",
			Partition: partition,
			Offset:    offset,
			Event:     event,
		}

		context.JSON(http.StatusCreated, response)
	}
}

func userEvent(producer sarama.SyncProducer, topic string) gin.HandlerFunc {
	return func(context *gin.Context) {

		var user dto.UserEvent
		if err := context.ShouldBindJSON(&user); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		event := dto.Event{
			ID:        fmt.Sprintf("user-%d-%s", user.UserID, user.Action),
			Type:      "user",
			Timestamp: time.Now().UTC(),
			Payload:   user,
		}

		msg, err := broker.NewProducerMessage(topic, event.Type, event)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// msg := broker.NewProducerMessage(topic, event.Type, event)
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := dto.EventResponse{
			Status:    "success",
			Partition: partition,
			Offset:    offset,
			Event:     event,
		}

		context.JSON(http.StatusCreated, response)
	}
}

func paymentEvent(producer sarama.SyncProducer, topic string) gin.HandlerFunc {
	return func(context *gin.Context) {

		var payment dto.PaymentEvent
		if err := context.ShouldBindJSON(&payment); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		event := dto.Event{
			ID:        fmt.Sprintf("payment-%d", payment.PaymentID),
			Type:      "payment",
			Timestamp: time.Now().UTC(),
			Payload:   payment,
		}

		msg, err := broker.NewProducerMessage(topic, event.Type, event)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := dto.EventResponse{
			Status:    "success",
			Partition: partition,
			Offset:    offset,
			Event:     event,
		}

		context.JSON(http.StatusCreated, response)
	}
}

func messageHandler(msg *sarama.ConsumerMessage) {
	log.Printf("Received message: topic=%s partition=%d offset=%d key=%s value=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
}
