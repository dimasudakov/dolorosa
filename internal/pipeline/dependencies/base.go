package dependencies

import (
	"context"
	"fmt"

	"dolorosa/internal/pipeline/contracts"
)

// BaseDependency - хранит подзависимости и позволяет задавать зависимостям определенные в Options характеристики
type BaseDependency[O contracts.OperationData] struct {
	name    string
	subDeps []contracts.Dependency[O]
	options Options
}

// NewBaseDependency - конструктор
func NewBaseDependency[O contracts.OperationData](name string, subDeps []contracts.Dependency[O], opts ...OptionFunc) BaseDependency[O] {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return BaseDependency[O]{
		name:    name,
		subDeps: subDeps,
		options: options,
	}
}

// Name - название зависимости
func (b *BaseDependency[O]) Name() string {
	return b.name
}

// SubDependencies - геттер подзависимостей
func (b *BaseDependency[O]) SubDependencies() []contracts.Dependency[O] {
	return b.subDeps
}

// ResolveByFunc - обвязка вокруг ресолвинга зависимости
func (b *BaseDependency[O]) ResolveByFunc(
	ctx context.Context,
	resolveFunc func(ctx context.Context) error,
) error {
	fmt.Printf("[BaseDependency] name: %s resolving...\n", b.name)

	err := resolveFunc(ctx)

	if err != nil {
		if !b.options.Optional {
			return fmt.Errorf("can't resolve dependency %s: %v", b.name, err)
		}
		fmt.Printf("can't resolve dependency: %s: %v\n", b.name, err)
	}

	return nil
}

// Options - доп параметры для загрузки зависимости
type Options struct {
	Optional bool
}

// OptionFunc - ф-ия для установки опций
type OptionFunc func(*Options)

// WithOptional - функция для установки Optional
func WithOptional() OptionFunc {
	return func(opts *Options) {
		opts.Optional = true
	}
}
