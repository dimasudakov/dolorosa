package dolorosa

import (
	"context"
	"dolorosa/internal/pipeline/operations/sbp/domain"
	"dolorosa/pkg/api/control"
	"github.com/samber/lo"
	"log"
)

func (s *OnlineControlService) CheckSBP(ctx context.Context, req *control.CheckSBPRequest) (*control.CheckSBPResponse, error) {
	log.Printf("Received request: %+v", req)

	decision := s.sbpChecker.Start(ctx, domain.Operation{
		OperationID:   req.OperationId,
		ClientID:      req.ClientId,
		Amount:        req.Amount,
		SenderPhone:   req.SenderPhone,
		ReceiverPhone: req.ReceiverPhone,
		ReceiverBic:   lo.FromPtr(req.ReceiverBic),
		ReceiverName:  lo.FromPtr(req.ReceiverName),
	})

	return &control.CheckSBPResponse{
		Decision:      control.Decision(decision.Decision),
		DeclineReason: &decision.Reason,
	}, nil
}
