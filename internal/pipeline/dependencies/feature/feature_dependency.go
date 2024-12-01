package feature

import (
	"context"

	"dolorosa/internal/pipeline/contracts"

	deps "dolorosa/internal/pipeline/dependencies"
)

type featureDependency[O contracts.OperationData] struct {
	deps.BaseDependency[O]
	features []FeatureRequest
}

func NewFeatureDependency[O contracts.OperationData](
	featureGroupName string,
	subDeps []contracts.Dependency[O],
	features []FeatureRequest,
	opts ...deps.OptionFunc,
) contracts.Dependency[O] {
	return &featureDependency[O]{
		BaseDependency: deps.NewBaseDependency(featureGroupName, subDeps, opts...),
		features:       features,
	}
}

func (f *featureDependency[O]) Resolve(ctx context.Context, state contracts.State[O]) error {
	return f.ResolveByFunc(ctx, f.resolve(state))
}

func (f *featureDependency[O]) resolve(state contracts.State[O]) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		_, err := deps.AssertAdapter[Adapter](state.GetOperation(), f.Name())
		if err != nil {
			return err
		}

		return nil
	}
}
