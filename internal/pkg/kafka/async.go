package producer

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/IBM/sarama"
)

// SendMessageAsync отправка сообщения асинхронно, но через синхронный консьюмер
func (s *AsyncProducer) SendMessageAsync(ctx context.Context, messages ...Message) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SendMessageAsync panic: %s\n", err)
			}
		}()
		s.wg.Add(1)
		defer s.wg.Done()

		kafkaMessages := transformToProducerMessages(ctx, messages...)

		for i := range kafkaMessages {
			_, _, err := s.producer.SendMessage(kafkaMessages[i])
			if err != nil {
				fmt.Printf("producer.SendMessageWithContext: %v\n", err)
			}
		}

	}()
}

func transformToProducerMessages(ctx context.Context, messages ...Message) []*sarama.ProducerMessage {
	directMessage := make([]sarama.ProducerMessage, len(messages))
	pointerMessage := make([]*sarama.ProducerMessage, 0, len(messages))

	marshaller := protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	for i := range messages {
		if messages[i] == nil {
			continue
		}

		value, err := marshaller.Marshal(messages[i].Value())
		if err != nil {
			fmt.Printf("marshaller.Marshal: %v", err)
			continue
		}

		directMessage[i] = sarama.ProducerMessage{
			Topic: messages[i].Topic(ctx),
			Key:   sarama.StringEncoder(messages[i].Key()),
			Value: sarama.ByteEncoder(value),
		}

		pointerMessage = append(pointerMessage, &directMessage[i])
	}

	return pointerMessage
}
