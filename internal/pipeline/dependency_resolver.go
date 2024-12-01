package pipeline

import (
	"context"
	"fmt"
	"sync"

	"dolorosa/internal/pipeline/contracts"

	"dolorosa/internal/pipeline/utils"
	"github.com/samber/lo"
)

type dependencyResolver[T contracts.OperationData] struct{}

func NewDependencyResolver[T contracts.OperationData]() contracts.DependencyResolver[T] {
	return dependencyResolver[T]{}
}

func (d dependencyResolver[T]) Resolve(ctx context.Context, rules []contracts.Rule[T], state contracts.State[T]) <-chan contracts.Rule[T] {
	resultRuleChan := make(chan contracts.Rule[T], len(rules))
	wg := sync.WaitGroup{}

	uniqueDeps := d.collectDependencies(ctx, rules)
	state.InitStatuses(lo.Map(uniqueDeps, func(dep contracts.Dependency[T], _ int) string {
		return dep.Name()
	}))
	utils.Go(ctx, func(ctx context.Context) {
		d.resolveDeps(ctx, uniqueDeps, state)
	})

	for _, rule := range rules {
		wg.Add(1)

		utils.Go(ctx, func(ctx context.Context) {
			defer wg.Done()
			skipRule := false

			for _, dep := range rule.GetDependencies(ctx) {
				state.WaitResolving(dep.Name())
				if state.GetError(dep.Name()) != nil {
					skipRule = true
					fmt.Printf("[DependencyResolver] skip rule %s, error: %s\n", rule.Name(), state.GetError(dep.Name()).Error())
					break
				}
			}

			if !skipRule {
				fmt.Printf("[DependencyResolver] rule %s deps resolved\n", rule.Name())
				resultRuleChan <- rule
			}
		})
	}

	utils.Go(ctx, func(ctx context.Context) {
		wg.Wait()
		close(resultRuleChan)
	})

	return resultRuleChan
}

func (d dependencyResolver[T]) collectDependencies(ctx context.Context, rules []contracts.Rule[T]) (result []contracts.Dependency[T]) {
	uniqueDeps := make(map[string]struct{})

	var dfs func(dep contracts.Dependency[T])

	dfs = func(dep contracts.Dependency[T]) {
		if _, ok := uniqueDeps[dep.Name()]; !ok {
			result = append(result, dep)
		}
		uniqueDeps[dep.Name()] = struct{}{}
		for _, subDep := range dep.SubDependencies() {
			dfs(subDep)
		}
	}

	for _, rule := range rules {
		for _, dep := range rule.GetDependencies(ctx) {
			dfs(dep)
		}
	}

	return result
}

func (d dependencyResolver[T]) resolveDeps(ctx context.Context, deps []contracts.Dependency[T], state contracts.State[T]) {
	for _, dep := range deps {

		utils.Go(ctx, func(ctx context.Context) {
			for _, subDep := range dep.SubDependencies() {
				state.WaitResolving(subDep.Name())
				if subDepErr := state.GetError(subDep.Name()); subDepErr != nil {
					state.MarkResolved(dep.Name(), fmt.Errorf("can't resolve dependency: %s, err: %v", dep.Name(), subDepErr))
					return
				}
			}
			err := dep.Resolve(ctx, state)
			state.MarkResolved(dep.Name(), err)
		})

	}
}
