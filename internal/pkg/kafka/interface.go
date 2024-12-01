package producer

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Message сообщение для отправки
type Message interface {
	// Topic название топика куда постить
	Topic(ctx context.Context) string
	// Key ключ, по которому сообщение распределят в партицию
	Key() string
	// Value значение
	Value() proto.Message
}
