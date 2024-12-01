package dependencies

import (
	"context"
	"testing"

	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBaseDependency(t *testing.T) {
	t.Parallel()

	type fields struct {
		name    string
		subDeps []contracts.Dependency[MockOperationData]
		options []OptionFunc
	}

	type args struct {
		ctx         context.Context
		resolveFunc func(ctx context.Context) error
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		expected expected
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				resolveFunc: func(ctx context.Context) error {
					return nil
				},
			},
			fields: fields{
				name: "base_dep",
				subDeps: []contracts.Dependency[MockOperationData]{
					mock.NewMockDependency[MockOperationData](gomock.NewController(t)),
					mock.NewMockDependency[MockOperationData](gomock.NewController(t)),
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				resolveFunc: func(ctx context.Context) error {
					return errors.New("some_err")
				},
			},
			fields: fields{
				name: "base_dep",
				subDeps: []contracts.Dependency[MockOperationData]{
					mock.NewMockDependency[MockOperationData](gomock.NewController(t)),
					mock.NewMockDependency[MockOperationData](gomock.NewController(t)),
				},
			},
			expected: expected{
				err: errors.New("can't resolve dependency base_dep: some_err"),
			},
		},
		{
			name: "error, but optional",
			args: args{
				ctx: context.Background(),
				resolveFunc: func(ctx context.Context) error {
					return errors.New("some_err")
				},
			},
			fields: fields{
				name: "base_dep",
				subDeps: []contracts.Dependency[MockOperationData]{
					mock.NewMockDependency[MockOperationData](gomock.NewController(t)),
					mock.NewMockDependency[MockOperationData](gomock.NewController(t)),
				},
				options: []OptionFunc{
					WithOptional(),
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			baseDep := NewBaseDependency[MockOperationData](tt.fields.name, tt.fields.subDeps, tt.fields.options...)

			assert.Equal(t, tt.fields.name, baseDep.Name())
			assert.Equal(t, tt.fields.subDeps, baseDep.SubDependencies())

			err := baseDep.ResolveByFunc(tt.args.ctx, tt.args.resolveFunc)
			if tt.expected.err != nil {
				assert.EqualError(t, err, tt.expected.err.Error())
			}
		})
	}
}

type MockOperationData struct {
	Receiver string
	Amount   int
}

type MockRuleDecision struct {
	RuleName string
	Decision string
}

type MockFinalDecision struct {
	FinalDecision string
	Info          string
}
