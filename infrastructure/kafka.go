package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type KafkaProducer struct {
	Conn *kafka.Conn
}

func NewKafka(ctx context.Context, topic string) (*KafkaProducer, error) {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	if kafkaHost == "" || kafkaPort == "" {
		return nil, errors.New("KAFKA_HOST and KAFKA_PORT environment variables are required")
	}
	address := fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)
	log.Printf("connecting to Kafka at %s", address)
	conn, err := kafka.DialLeader(ctx, "tcp", address, topic, 0)
	if err != nil {
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		Conn: conn,
	}, nil
}
