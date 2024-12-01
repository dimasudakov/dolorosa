package sbp

import (
	"context"
	"dolorosa/internal/pipeline/common"
	"dolorosa/internal/pkg/notifier"
	"fmt"

	"dolorosa/internal/pipeline/contracts"

	"dolorosa/internal/pipeline/utils"
)

type b2cDecisionResolver struct {
	notifier notifier.Notifier
}

func newSbpDecisionResolver(fields CheckerFields) contracts.DecisionResolver {
	return b2cDecisionResolver{
		notifier: fields.Notifier,
	}
}

func (b b2cDecisionResolver) Resolve(ctx context.Context, ch <-chan common.RuleDecision) (result common.FinalDecision) {

	resChan := make(chan common.FinalDecision)

	utils.Go(ctx, func(ctx context.Context) {
		defer close(resChan)

		for ruleDecision := range ch {
			if ruleDecision.Decision == common.Decline {

				if ruleDecision.AlertInfo != nil {
					utils.Go(ctx, func(ctx context.Context) {
						b.sendAlert(ctx, ruleDecision)
					})
				}

				select {
				case resChan <- common.FinalDecision{Decision: common.Decline}:
				default:
				}
			}
		}

		// Если ни одно правило не выставило shouldDecline, отправляем положительный результат
		select {
		case resChan <- common.FinalDecision{Decision: common.Ok}:
		default:
		}
	})

	defer func() {
		fmt.Printf("[Sbp Resolver] Resolve complete with decision: %+v\n", result)
	}()

	return <-resChan
}

func (b b2cDecisionResolver) sendAlert(ctx context.Context, decision common.RuleDecision) {
	info := decision.AlertInfo
	err := b.notifier.SendNotification(ctx, notifier.Notification{
		Text:           info.Msg,
		AlertName:      info.Name,
		ClientID:       info.ClientID,
		Amount:         info.Amount,
		IdempotencyKey: &info.IdempotencyKey,
		RuleName:       decision.RuleName,
		OperationID:    info.OperationID,
	})
	if err != nil {
		fmt.Printf("can't send alert, error: %s", err.Error())
	}
}
