package models

import (
	"context"
	contracts "dolorosa/pkg/api/kafka"
	"google.golang.org/protobuf/proto"
)

// DecisionLog - обертка над логом для отправки в кафку
type DecisionLog struct {
	msg *contracts.DecisionLog
}

// NewDecisionLog - конструктор
func NewDecisionLog(msg *contracts.DecisionLog) *DecisionLog {
	return &DecisionLog{msg: msg}
}

// Topic - топик, в который должен быть отправлен лог
func (s *DecisionLog) Topic(_ context.Context) string {
	return "logs"
}

// Key - ключ для распределения по партициям
func (s *DecisionLog) Key() string {
	if s == nil {
		return ""
	}
	return s.msg.GetOperationId()
}

// Value - сообщение, которое будет передано в кафку
func (s *DecisionLog) Value() proto.Message {
	if s == nil {
		return nil
	}
	return s.msg
}
