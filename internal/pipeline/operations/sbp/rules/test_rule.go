package rules

import (
	"context"
	"dolorosa/internal/pipeline/common"
	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/operations/sbp/domain"
	"fmt"
)

type testRule struct {
	dependencies []contracts.Dependency[domain.Operation]
}

func NewTestRule(deps []contracts.Dependency[domain.Operation]) contracts.Rule[domain.Operation] {
	return testRule{
		dependencies: deps,
	}
}

func (t testRule) Name() string {
	return "test_rule"
}

func (t testRule) ShouldRun(_ context.Context, _ domain.Operation) bool {
	return true
}

func (t testRule) Check(_ context.Context, state contracts.State[domain.Operation]) (decision common.RuleDecision) {
	exception := state.GetExceptions(t.Name())
	if exception.Found {
		fmt.Printf("[Test rule] exception for rule: %s + clientID: %s found, skip rule", t.Name(), state.GetOperation().ClientID)
		return common.RuleDecision{
			Decision: common.Ok,
		}
	}

	if state.GetOperation().Amount > 5000 {
		return common.RuleDecision{
			Decision: common.Decline,
			Reason:   t.Name(),
		}
	}

	return common.RuleDecision{
		Decision: common.Ok,
	}
}

func (t testRule) GetDependencies(_ context.Context) []contracts.Dependency[domain.Operation] {
	return t.dependencies
}
