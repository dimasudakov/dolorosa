package interceptors

import (
	"context"
	"dolorosa/internal/pkg/utils"
	"dolorosa/pkg/api/control"
	contracts "dolorosa/pkg/api/kafka"
	"fmt"
	"google.golang.org/grpc"
)

type LogSender interface {
	SendLog(ctx context.Context, log *contracts.DecisionLog) error
}

type AuditLogsInterceptor struct {
	logSender LogSender
}

func NewAuditLogsInterceptor(logSender LogSender) *AuditLogsInterceptor {
	return &AuditLogsInterceptor{
		logSender: logSender,
	}
}

func (i *AuditLogsInterceptor) Unary() grpc.UnaryServerInterceptor {
	return i.handle
}

func (i *AuditLogsInterceptor) handle(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	next grpc.UnaryHandler,
) (resp interface{}, respErr error) {
	resp, respErr = next(ctx, req)

	decisionLog, err := i.buildLog(ctx, req, resp)
	if err != nil {
		return resp, respErr
	}

	_ = i.logSender.SendLog(ctx, decisionLog)

	return resp, respErr
}

func (i *AuditLogsInterceptor) buildLog(ctx context.Context, req interface{}, resp interface{}) (*contracts.DecisionLog, error) {
	log := &contracts.DecisionLog{}

	err := i.fillLogFromRequest(log, req)
	if err != nil {
		return nil, err
	}

	err = i.fillLogFromResponse(log, resp)
	if err != nil {
		return nil, err
	}

	log.TraceId = utils.GetTraceID(ctx)

	return log, nil
}

func (i *AuditLogsInterceptor) fillLogFromRequest(log *contracts.DecisionLog, req interface{}) error {
	switch r := req.(type) {
	case *control.CheckSBPRequest:
		log.OperationId = r.OperationId
		log.ClientId = r.ClientId
		log.Amount = r.Amount
	default:
		return fmt.Errorf("unknown request type %T", r)
	}

	return nil
}

func (i *AuditLogsInterceptor) fillLogFromResponse(log *contracts.DecisionLog, resp interface{}) error {
	switch r := resp.(type) {
	case *control.CheckSBPResponse:
		log.Decision = r.Decision.String()
		if r.DeclineReason != nil {
			log.DeclineReason = r.DeclineReason
		}
	default:
		return fmt.Errorf("unknown response type %T", r)
	}
	return nil
}
