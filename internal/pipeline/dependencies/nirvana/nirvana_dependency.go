package nirvana

import (
	"context"
	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pkg/nirvana_helper"

	deps "dolorosa/internal/pipeline/dependencies"
)

const nirvanaDependencyName = "nirvana_dependency"

type ExceptionChecker interface {
	CheckException(ctx context.Context, name string, attributes nirvana_helper.ExceptionAttributes) (bool, error)
}

type nirvanaDependency[O contracts.OperationData] struct {
	deps.BaseDependency[O]
	exceptions       []Exception
	exceptionChecker ExceptionChecker
}

func NewNirvanaDependency[O contracts.OperationData](exceptions []Exception, checker ExceptionChecker, opts ...deps.OptionFunc) contracts.Dependency[O] {
	return &nirvanaDependency[O]{
		BaseDependency:   deps.NewBaseDependency[O](nirvanaDependencyName, []contracts.Dependency[O]{}, opts...),
		exceptions:       exceptions,
		exceptionChecker: checker,
	}
}

func (n *nirvanaDependency[O]) Resolve(ctx context.Context, state contracts.State[O]) error {
	return n.ResolveByFunc(ctx, n.resolve(state))
}

func (n *nirvanaDependency[O]) resolve(state contracts.State[O]) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		adapter, err := deps.AssertAdapter[Adapter](state.GetOperation(), n.Name())
		if err != nil {
			return err
		}

		for _, exception := range n.exceptions {
			found, fErr := n.exceptionChecker.CheckException(
				ctx,
				exception.Name,
				nirvana_helper.ExceptionAttributes{
					ClientID: adapter.GetClientID(),
				},
			)
			if fErr != nil {
				continue
			}
			state.SetExceptions(exception.Name, contracts.ExceptionInfo{
				Found: found,
			})
		}

		return nil
	}
}
