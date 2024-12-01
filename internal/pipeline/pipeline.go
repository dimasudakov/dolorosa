package pipeline

import (
	"context"

	"github.com/samber/lo"

	"dolorosa/internal/pipeline/common"
	"dolorosa/internal/pipeline/contracts"
)

type pipeline[O contracts.OperationData] struct {
	ruleRegistry       contracts.RuleRegistry[O]
	dependencyResolver contracts.DependencyResolver[O]
	ruleExecutor       contracts.RuleExecutor[O]
	decisionResolver   contracts.DecisionResolver
}

func NewPipeline[O contracts.OperationData](
	ruleRegistry contracts.RuleRegistry[O],
	dependencyResolver contracts.DependencyResolver[O],
	ruleExecutor contracts.RuleExecutor[O],
	decisionResolver contracts.DecisionResolver,
) contracts.Pipeline[O] {
	return &pipeline[O]{
		ruleRegistry:       ruleRegistry,
		dependencyResolver: dependencyResolver,
		ruleExecutor:       ruleExecutor,
		decisionResolver:   decisionResolver,
	}
}

func (r pipeline[O]) Start(ctx context.Context, operation O) common.FinalDecision {
	rules := r.ruleRegistry.GetRules()

	rules = lo.Filter(rules, func(rule contracts.Rule[O], _ int) bool {
		return rule.ShouldRun(ctx, operation)
	})

	state := contracts.NewState(operation)

	depsResolvedCh := r.dependencyResolver.Resolve(ctx, rules, state)

	rulesResultCh := r.ruleExecutor.ExecuteRules(ctx, depsResolvedCh, state, len(rules))

	return r.decisionResolver.Resolve(ctx, rulesResultCh)
}
