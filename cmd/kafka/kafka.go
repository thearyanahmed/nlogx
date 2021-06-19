package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

func ReaderConfig(brokers []string,kafkaTopic ,kafkaClientId string) *kafka.ReaderConfig {
	// make a new reader that consumes from topic-A
	return &kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientId,
		Topic:           kafkaTopic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}
}

func Consume(config kafka.ReaderConfig) /*(kafka.Message,error)*/ {
	reader := kafka.NewReader(config)

	defer func(reader *kafka.Reader) {
		_ = reader.Close()
	}(reader)

	for {
		m, err := reader.ReadMessage(context.Background())
		fmt.Printf("tpc: %v | prtn: %v | ofst: %v |: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))

		if err != nil {
			fmt.Printf("error not nil, %v\n",err.Error())
		}
		//return m, err
	}
}

type WriterInterface interface {
	Writer(brokers []string,topic string) *kafka.Writer
	Produce(w *kafka.Writer, ctx context.Context, key, msg string) error
}