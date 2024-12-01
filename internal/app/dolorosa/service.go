package dolorosa

import (
	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/operations/sbp/domain"
	"dolorosa/pkg/api/control"
)

type OnlineControlService struct {
	sbpChecker contracts.Pipeline[domain.Operation]

	control.UnimplementedOnlineControlServer
}

func NewOnlineControlService(sbpChecker contracts.Pipeline[domain.Operation]) *OnlineControlService {
	return &OnlineControlService{
		sbpChecker: sbpChecker,
	}
}
