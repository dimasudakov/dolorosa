package dependencies

import (
	"fmt"
)

// AssertAdapter - приводит state.Operation к нужному адаптеру для дальнейшего разрешения зависимости
func AssertAdapter[T any](operation any, dependencyName string) (T, error) {
	adapter, ok := operation.(T)
	if !ok {
		return *new(T), fmt.Errorf("operation does not implement required for dependency: %s, adapter interface: %T", dependencyName, operation)
	}
	return adapter, nil
}

func MaybeAdapter[T any](operation any) (T, bool) {
	adapter, ok := operation.(T)
	if !ok {
		return *new(T), false
	}
	return adapter, true
}
