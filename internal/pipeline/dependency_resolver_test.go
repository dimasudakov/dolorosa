package pipeline

import (
	"context"
	"sync"
	"testing"
	"time"

	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDependencyResolver_Resolve(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		rules []contracts.Rule[MockOperationData]
		state contracts.State[MockOperationData]
	}

	type expected struct {
		ruleNames []string
	}

	tests := []struct {
		name     string
		args     func(ctrl *gomock.Controller) args
		expected expected
	}{
		{
			name: "success",
			args: func(ctrl *gomock.Controller) args {

				state := contracts.NewState[MockOperationData](MockOperationData{Amount: 5000_00, Receiver: "dima"})

				rule1 := mock.NewMockRule[MockOperationData](ctrl)
				rule2 := mock.NewMockRule[MockOperationData](ctrl)
				rule3 := mock.NewMockRule[MockOperationData](ctrl)
				rule4 := mock.NewMockRule[MockOperationData](ctrl)

				rule1.EXPECT().Name().AnyTimes().Return("rule_1")
				rule2.EXPECT().Name().AnyTimes().Return("rule_2")
				rule3.EXPECT().Name().AnyTimes().Return("rule_3")
				rule4.EXPECT().Name().AnyTimes().Return("rule_4")

				rule1Dep1 := testDependency{
					name:              "rule_1_dep_1",
					resolvingDuration: time.Millisecond * 50,
					resolvingError:    nil,
					mu:                &sync.Mutex{},
					subDeps: []testDependency{
						{
							name:              "rule_1_dep_1_1",
							resolvingDuration: time.Millisecond * 10,
							resolvingError:    nil,
							mu:                &sync.Mutex{},
							subDeps:           []testDependency{},
						},
						{
							name:              "rule_1_dep_1_2",
							resolvingDuration: time.Millisecond * 100,
							resolvingError:    nil,
							mu:                &sync.Mutex{},
							subDeps: []testDependency{
								{
									name:              "rule_1_dep_1_2_1",
									resolvingDuration: time.Millisecond * 70,
									resolvingError:    nil,
									mu:                &sync.Mutex{},
									subDeps:           []testDependency{},
								},
							},
						},
					},
				}
				rule1Dep2 := testDependency{
					name:              "rule_1_dep_2",
					resolvingDuration: time.Millisecond * 100,
					resolvingError:    nil,
					mu:                &sync.Mutex{},
					subDeps:           []testDependency{},
				}

				rule3Dep1 := testDependency{
					name:              "rule_3_dep_1",
					resolvingDuration: time.Millisecond * 300,
					resolvingError:    nil,
					mu:                &sync.Mutex{},
					subDeps:           []testDependency{},
				}

				rule4Dep1 := testDependency{
					name:              "rule_4_dep_1",
					resolvingDuration: time.Millisecond * 100,
					resolvingError:    nil,
					mu:                &sync.Mutex{},
					subDeps: []testDependency{
						{
							name:              "rule_4_dep_1_1",
							resolvingDuration: time.Millisecond * 10,
							resolvingError:    nil,
							mu:                &sync.Mutex{},
							subDeps: []testDependency{
								{
									name:              "rule_4_dep_1_1_1",
									resolvingDuration: time.Millisecond * 70,
									resolvingError:    errors.New("invalid argument request"),
									mu:                &sync.Mutex{},
									subDeps: []testDependency{
										{
											name:              "rule_4_dep_1_1_1_1",
											resolvingDuration: time.Millisecond * 10,
											resolvingError:    nil,
											mu:                &sync.Mutex{},
											subDeps:           []testDependency{},
										},
									},
								},
							},
						},
					},
				}

				rule1.EXPECT().GetDependencies(gomock.Any()).Times(3).Return([]contracts.Dependency[MockOperationData]{rule1Dep1, rule1Dep2})
				rule2.EXPECT().GetDependencies(gomock.Any()).Times(3).Return([]contracts.Dependency[MockOperationData]{})
				rule3.EXPECT().GetDependencies(gomock.Any()).Times(3).Return([]contracts.Dependency[MockOperationData]{rule3Dep1})
				rule4.EXPECT().GetDependencies(gomock.Any()).Times(2).Return([]contracts.Dependency[MockOperationData]{rule4Dep1})

				return args{
					ctx: context.Background(),
					rules: []contracts.Rule[MockOperationData]{
						// ~220ms
						rule1,
						// ~0ms
						rule2,
						// ~300ms
						rule3,
						// ~80ms but error
						rule4,
					},
					state: state,
				}
			},
			expected: expected{
				ruleNames: []string{"rule_2", "rule_1", "rule_3"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			depsState = newResolvedState()
			ctrl := gomock.NewController(t)
			a := tt.args(ctrl)
			depResolver := NewDependencyResolver[MockOperationData]()

			resultCh := depResolver.Resolve(a.ctx, a.rules, a.state)

			for rule := range resultCh {
				assert.Equal(t, tt.expected.ruleNames[0], rule.Name())

				for _, dep := range rule.GetDependencies(a.ctx) {
					assert.True(t, depsState.isResolved(dep.Name()))
				}

				tt.expected.ruleNames = tt.expected.ruleNames[1:]
			}
		})
	}
}

// --------------------------------- helpers ---------------------------------

// resolvedState используем в тестах, чтоб проверить корректность состояния зависимостей перед их запуском
type resolvedState struct {
	r  map[string]bool
	mu *sync.RWMutex
}

func newResolvedState() resolvedState {
	return resolvedState{
		r:  make(map[string]bool),
		mu: &sync.RWMutex{},
	}
}

func (rs *resolvedState) setResolved(name string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.r[name] = true
}

func (rs *resolvedState) isResolved(name string) bool {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	return rs.r[name]
}

var depsState resolvedState

// testDependency - тестовая зависимость
type testDependency struct {
	name              string
	resolvingDuration time.Duration
	resolvingError    error
	subDeps           []testDependency

	mu *sync.Mutex
}

func (t testDependency) Name() string {
	return t.name
}

func (t testDependency) Resolve(_ context.Context, _ contracts.State[MockOperationData]) error {
	if depsState.isResolved(t.name) {
		panic("dependency should not be resolved more than once")
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, dep := range t.subDeps {
		if !depsState.isResolved(dep.name) {
			panic("all dependency's subdeps should be resolved")
		}
	}
	depsState.setResolved(t.name)

	time.Sleep(t.resolvingDuration)
	return t.resolvingError
}

func (t testDependency) SubDependencies() []contracts.Dependency[MockOperationData] {
	res := make([]contracts.Dependency[MockOperationData], 0, len(t.subDeps))
	for _, dep := range t.subDeps {
		res = append(res, dep)
	}
	return res
}
