package services

import "github.com/segmentio/kafka-go"

type KafkaProducer struct {
	conn *kafka.Conn
}

func NewKafkaProducer(conn *kafka.Conn) *KafkaProducer {
	return &KafkaProducer{
		conn: conn,
	}
}

func (kp *KafkaProducer) WriteMsg(msg string) error {
	_, err := kp.conn.WriteMessages(
		kafka.Message{
			Value: []byte(msg),
		})
	if err != nil {
		return err
	}
	return nil
}
