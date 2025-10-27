package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ploezy/ecommerce-platform/order-service/pkg/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	cfg    *config.Config
}

var producer *Producer

func NewProducer(cfg *config.Config) (*Producer, error) {
	// Parse Kafka brokers
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	// Create Kafka writer
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}

	producer = &Producer{
		writer: writer,
		cfg:    cfg,
	}

	log.Println(" Kafka producer initialized successfully")
	return producer, nil
}

func GetProducer() *Producer {
	return producer
}

// PublishEvent publishes an event to a specific topic
func (p *Producer) PublishEvent(topic string, key string, data interface{}) error {
	// Convert data to JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create message
	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}

	// Send message
	ctx := context.Background()
	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf(" Event published to topic [%s] with key [%s]", topic, key)
	return nil
}

// Close closes the Kafka writer
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}