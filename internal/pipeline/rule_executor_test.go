package pipeline

import (
	"context"
	"slices"
	"strings"
	"testing"

	"dolorosa/internal/pipeline/common"

	"dolorosa/internal/pipeline/contracts"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dolorosa/internal/pipeline/mock"
)

func TestRuleExecutor_ExecuteRules(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		rules    []contracts.Rule[MockOperationData]
		state    contracts.State[MockOperationData]
		rulesCnt int
	}

	type expected struct {
		results []common.RuleDecision
	}

	tests := []struct {
		name        string
		args        func(ctrl *gomock.Controller) args
		ruleChannel func(rules []contracts.Rule[MockOperationData]) <-chan contracts.Rule[MockOperationData]
		expected    expected
		expectErr   bool
	}{
		{
			name: "success - all rules execute successfully",
			args: func(ctrl *gomock.Controller) args {

				state := mock.NewMockState[MockOperationData](ctrl)

				rule1 := mock.NewMockRule[MockOperationData](ctrl)
				rule2 := mock.NewMockRule[MockOperationData](ctrl)
				rule3 := mock.NewMockRule[MockOperationData](ctrl)

				rule1.EXPECT().Check(gomock.Any(), state).Times(1).Return(common.RuleDecision{
					RuleName: "rule_1",
					Decision: common.Ok,
				})
				rule2.EXPECT().Check(gomock.Any(), state).Times(1).Return(common.RuleDecision{
					RuleName: "rule_2",
					Decision: common.Decline,
				})
				rule3.EXPECT().Check(gomock.Any(), state).Times(1).Return(common.RuleDecision{
					RuleName: "rule_3",
					Decision: common.Ok,
				})
				rule1.EXPECT().Name().AnyTimes().Return("rule_1_name")
				rule2.EXPECT().Name().AnyTimes().Return("rule_2_name")
				rule3.EXPECT().Name().AnyTimes().Return("rule_3_name")

				return args{
					ctx: context.Background(),
					rules: []contracts.Rule[MockOperationData]{
						rule1,
						rule2,
						rule3,
					},
					state:    state,
					rulesCnt: 3,
				}
			},
			ruleChannel: func(rules []contracts.Rule[MockOperationData]) <-chan contracts.Rule[MockOperationData] {
				ch := make(chan contracts.Rule[MockOperationData])

				go func() {
					defer close(ch)
					ch <- rules[2]
					ch <- rules[0]
					ch <- rules[1]
				}()

				return ch
			},
			expected: expected{
				results: []common.RuleDecision{
					{RuleName: "rule_1", Decision: common.Ok},
					{RuleName: "rule_2", Decision: common.Decline},
					{RuleName: "rule_3", Decision: common.Ok},
				},
			},
			expectErr: false,
		},
		{
			name: "panic in rule",
			args: func(ctrl *gomock.Controller) args {

				state := mock.NewMockState[MockOperationData](ctrl)

				rule1 := mock.NewMockRule[MockOperationData](ctrl)
				rule2 := mock.NewMockRule[MockOperationData](ctrl)
				rule3 := mock.NewMockRule[MockOperationData](ctrl)

				rule1.EXPECT().Check(gomock.Any(), state).Times(1).Return(common.RuleDecision{
					RuleName: "rule_1",
					Decision: common.Ok,
				})
				rule2.EXPECT().Check(gomock.Any(), state).Times(1).Return(common.RuleDecision{
					RuleName: "rule_2",
					Decision: common.Decline,
				})
				rule3.EXPECT().Check(gomock.Any(), state).Times(1).Do(
					func(ctx context.Context, state contracts.State[MockOperationData]) {
						panic("panic in rule_3")
					},
				)
				rule1.EXPECT().Name().AnyTimes().Return("rule_1_name")
				rule2.EXPECT().Name().AnyTimes().Return("rule_2_name")
				rule3.EXPECT().Name().AnyTimes().Return("rule_3_name")

				return args{
					ctx: context.Background(),
					rules: []contracts.Rule[MockOperationData]{
						rule1,
						rule2,
						rule3,
					},
					state:    state,
					rulesCnt: 3,
				}
			},
			ruleChannel: func(rules []contracts.Rule[MockOperationData]) <-chan contracts.Rule[MockOperationData] {
				ch := make(chan contracts.Rule[MockOperationData])

				go func() {
					defer close(ch)
					ch <- rules[1]
					ch <- rules[2]
					ch <- rules[0]
				}()

				return ch
			},
			expected: expected{
				results: []common.RuleDecision{
					{RuleName: "rule_1", Decision: common.Ok},
					{RuleName: "rule_2", Decision: common.Decline},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			a := tt.args(ctrl)
			executor := NewRuleExecutor[MockOperationData]()
			rulesCh := tt.ruleChannel(a.rules)

			resultsCh := executor.ExecuteRules(a.ctx, rulesCh, a.state, a.rulesCnt)

			results := make([]common.RuleDecision, 0)
			for result := range resultsCh {
				results = append(results, result)
			}
			slices.SortFunc(results, func(a, b common.RuleDecision) int {
				return strings.Compare(a.RuleName, b.RuleName)
			})
			assert.Equal(t, tt.expected.results, results)
		})
	}
}

// ----------------------------------------- Mock data types -----------------------------------------

type MockOperationData struct {
	Receiver string
	Amount   int
}
