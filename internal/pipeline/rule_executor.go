package pipeline

import (
	"context"
	"fmt"
	"sync"

	"dolorosa/internal/pipeline/common"

	"dolorosa/internal/pipeline/contracts"

	"dolorosa/internal/pipeline/utils"
)

type ruleExecutor[O contracts.OperationData] struct{}

func NewRuleExecutor[O contracts.OperationData]() contracts.RuleExecutor[O] {
	return ruleExecutor[O]{}
}

func (r ruleExecutor[O]) ExecuteRules(
	ctx context.Context,
	rules <-chan contracts.Rule[O],
	state contracts.State[O],
	rulesCnt int,
) <-chan common.RuleDecision {
	rulesResultCh := make(chan common.RuleDecision, rulesCnt)

	wg := sync.WaitGroup{}
	for rule := range rules {
		wg.Add(1)

		utils.Go(ctx, func(ctx context.Context) {
			defer wg.Done()

			fmt.Printf("[Engine] rule %s start\n", rule.Name())

			rulesResultCh <- rule.Check(ctx, state)

			fmt.Printf("[Engine] rule %s done\n", rule.Name())
		})
	}

	utils.Go(ctx, func(ctx context.Context) {
		wg.Wait()
		close(rulesResultCh)
	})

	return rulesResultCh
}
