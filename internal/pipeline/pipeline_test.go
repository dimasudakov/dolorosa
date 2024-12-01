package pipeline

import (
	"context"
	"testing"

	"dolorosa/internal/pipeline/common"

	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPipeline_Start(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		operation MockOperationData
	}

	type mocks struct {
		ruleRegistry       contracts.RuleRegistry[MockOperationData]
		dependencyResolver contracts.DependencyResolver[MockOperationData]
		ruleExecutor       contracts.RuleExecutor[MockOperationData]
		decisionResolver   contracts.DecisionResolver
	}

	tests := []struct {
		name     string
		args     args
		mocks    func(ctrl *gomock.Controller, a args) mocks
		expected common.FinalDecision
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				operation: MockOperationData{
					Receiver: "dima",
					Amount:   5000_00,
				},
			},
			mocks: func(ctrl *gomock.Controller, a args) mocks {
				ruleRegistryMock := mock.NewMockRuleRegistry[MockOperationData](ctrl)
				depResolverMock := mock.NewMockDependencyResolver[MockOperationData](ctrl)
				ruleExecutorMock := mock.NewMockRuleExecutor[MockOperationData](ctrl)
				decisionResolverMock := mock.NewMockDecisionResolver(ctrl)

				rule1 := mock.NewMockRule[MockOperationData](ctrl)
				rule1.EXPECT().ShouldRun(gomock.Any(), a.operation).Return(false)

				rule2 := mock.NewMockRule[MockOperationData](ctrl)
				rule2.EXPECT().ShouldRun(gomock.Any(), a.operation).Return(true)

				state := contracts.NewState(a.operation)

				ruleRegistryMock.EXPECT().
					GetRules().
					Times(1).
					Return([]contracts.Rule[MockOperationData]{rule1, rule2})

				depsResolvedCh := make(chan contracts.Rule[MockOperationData])
				depResolverMock.EXPECT().
					Resolve(gomock.Any(), []contracts.Rule[MockOperationData]{rule2}, state).
					Times(1).
					Return(depsResolvedCh)

				rulesResultCh := make(chan common.RuleDecision)
				ruleExecutorMock.EXPECT().
					ExecuteRules(gomock.Any(), depsResolvedCh, state, 1).
					Times(1).
					Return(rulesResultCh)

				decisionResolverMock.EXPECT().
					Resolve(gomock.Any(), rulesResultCh).
					Times(1).
					Return(common.FinalDecision{
						Decision: common.Ok,
						Reason:   "some_reason",
					})

				return mocks{
					ruleRegistry:       ruleRegistryMock,
					dependencyResolver: depResolverMock,
					ruleExecutor:       ruleExecutorMock,
					decisionResolver:   decisionResolverMock,
				}
			},
			expected: common.FinalDecision{
				Decision: common.Ok,
				Reason:   "some_reason",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			m := tt.mocks(ctrl, tt.args)
			pipeline := NewPipeline(
				m.ruleRegistry,
				m.dependencyResolver,
				m.ruleExecutor,
				m.decisionResolver,
			)

			result := pipeline.Start(tt.args.ctx, tt.args.operation)

			assert.Equal(t, tt.expected, result)
		})
	}
}
