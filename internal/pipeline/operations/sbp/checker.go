package sbp

import (
	"dolorosa/internal/pipeline"
	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/dependencies/nirvana"
	"dolorosa/internal/pipeline/operations/sbp/domain"
	"dolorosa/internal/pkg/notifier"
)

type CheckerFields struct {
	Notifier         notifier.Notifier
	ExceptionChecker nirvana.ExceptionChecker
}

func NewCheckerSbp(fields CheckerFields) contracts.Pipeline[domain.Operation] {
	return pipeline.NewPipeline(
		newSbpRuleRegistry(fields),
		newSbpDependencyResolver(),
		newSbpRuleExecutor(),
		newSbpDecisionResolver(fields),
	)
}

func newSbpDependencyResolver() contracts.DependencyResolver[domain.Operation] {
	return pipeline.NewDependencyResolver[domain.Operation]()
}

func newSbpRuleExecutor() contracts.RuleExecutor[domain.Operation] {
	return pipeline.NewRuleExecutor[domain.Operation]()
}
