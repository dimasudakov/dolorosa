package logs_sender

import (
	"context"
	producer "dolorosa/internal/pkg/kafka"
	"dolorosa/internal/pkg/kafka/models"
	contracts "dolorosa/pkg/api/kafka"
	"encoding/json"
	"fmt"
)

type Producer interface {
	SendMessageAsync(ctx context.Context, messages ...producer.Message)
}

type logsSender struct {
	producer Producer
}

func NewLogsSender(producer producer.Producer) *logsSender {
	return &logsSender{
		producer: producer,
	}
}

func (ls *logsSender) SendLog(ctx context.Context, log *contracts.DecisionLog) error {
	jsonLog, _ := json.Marshal(log)

	fmt.Printf("[DecisionLog] Sending log: %s\n", string(jsonLog))

	ls.producer.SendMessageAsync(ctx, models.NewDecisionLog(log))

	return nil
}
