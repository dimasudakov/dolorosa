package contracts

import (
	"context"

	"dolorosa/internal/pipeline/common"
)

//go:generate mockgen -destination=../mock/mock_pipeline.go -package=mock -source=interfaces.go

// Dependency - общий интерфейс для всех зависимостей.
type Dependency[O OperationData] interface {
	Name() string
	Resolve(ctx context.Context, state State[O]) error
	SubDependencies() []Dependency[O]
}

// Rule - интерфейс для правил
type Rule[O OperationData] interface {
	Name() string
	// ShouldRun - решает, стоит ли запускать правило
	ShouldRun(ctx context.Context, operation O) bool
	// Check - запуск правила
	Check(ctx context.Context, state State[O]) (decision common.RuleDecision)
	// GetDependencies - возвращает все зависимости правила
	GetDependencies(ctx context.Context) []Dependency[O]
}

// RuleRegistry - интерфейс структуры, содержащей все правила для определенного типа операций
type RuleRegistry[O OperationData] interface {
	GetRules() []Rule[O]
}

// DependencyResolver - интерфейс сущности, которая занимается подгрузкой зависимостей для правил
type DependencyResolver[O OperationData] interface {
	Resolve(ctx context.Context, rules []Rule[O], state State[O]) <-chan Rule[O]
}

// RuleExecutor - интерфейс сущности, которая занимается исполнением правил
type RuleExecutor[O OperationData] interface {
	ExecuteRules(ctx context.Context, rules <-chan Rule[O], state State[O], rulesCnt int) <-chan common.RuleDecision
}

// DecisionResolver - интерфейс сущности, которая занимается принятием итогового решения на основе отработки правил
type DecisionResolver interface {
	// Resolve - принять решение по операции
	Resolve(ctx context.Context, ch <-chan common.RuleDecision) common.FinalDecision
}

// Pipeline - интерфейс пайплайна для запуска правил
type Pipeline[O OperationData] interface {
	// Start - запуск пайплайна проверок
	Start(ctx context.Context, operation O) common.FinalDecision
}

// State - интерфейс состояния работы пайплайна
type State[T OperationData] interface {
	GetOperation() T

	/* Состояние загрузки зависимостей */
	MarkResolved(name string, err error)
	GetStatus(name string) *ResolvingStatus
	GetError(name string) error
	WaitResolving(name string)
	InitStatuses(depNames []string)

	/* Зависимости */
	SetFeatures(name FeatureGroupName, features []FeatureInfo)
	GetFeatures(name FeatureGroupName) []FeatureInfo
	SetExceptions(name string, exceptions ExceptionInfo)
	GetExceptions(name string) ExceptionInfo
}
