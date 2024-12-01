package producer

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

//go:generate mockgen -destination=../mock/kafka/producer.go gitlab.ozon.ru/ob/sauron-ng/internal/pkg/kafka Producer

// Producer - интерфейс для отправки сообщений в кафку асинхронно
type Producer interface {
	SendMessageAsync(ctx context.Context, messages ...Message)
	CloseAndWait()
}

// AsyncProducer - реализация Producer для отправки в кафку
type AsyncProducer struct {
	producer sarama.SyncProducer
	wg       sync.WaitGroup
}

// New - конструктор
func New(ctx context.Context) (Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(
		[]string{"localhost:9091", "localhost:9092"},
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf(`sarama.NewSyncProducer: %w`, err)
	}

	return &AsyncProducer{
		producer: producer,
	}, nil
}
