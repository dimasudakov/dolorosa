package contracts

import (
	"sync"
)

type PipelineState[T OperationData] struct {
	Operation T
	DependencyState
}

func (s *PipelineState[T]) GetOperation() T {
	return s.Operation
}

func NewState[T OperationData](operation T) *PipelineState[T] {
	return &PipelineState[T]{
		Operation:       operation,
		DependencyState: NewDependencyState(),
	}
}

type DependencyState struct {
	features   map[FeatureGroupName][]FeatureInfo
	exceptions map[string]ExceptionInfo

	// key - dependency name, value - resolving status info
	statuses map[string]*ResolvingStatus

	mu sync.RWMutex
}

func NewDependencyState() DependencyState {
	return DependencyState{
		features:   make(map[FeatureGroupName][]FeatureInfo),
		exceptions: make(map[string]ExceptionInfo),

		statuses: make(map[string]*ResolvingStatus),
	}
}

// -------------------------------------- Dependency resolving methods --------------------------------------

func (d *DependencyState) InitStatuses(depNames []string) {
	statuses := make(map[string]*ResolvingStatus)
	for _, name := range depNames {
		statuses[name] = &ResolvingStatus{
			done: make(chan struct{}),
		}
	}
	d.statuses = statuses
}

type ResolvingStatus struct {
	err  error
	done chan struct{}
	once sync.Once
}

func (d *DependencyState) MarkResolved(name string, err error) {
	status := d.statuses[name]
	status.once.Do(func() {
		status.err = err
		close(status.done)
	})
}

func (d *DependencyState) GetStatus(name string) *ResolvingStatus {
	return d.statuses[name]
}

func (d *DependencyState) GetError(name string) error {
	return d.GetStatus(name).err
}

func (d *DependencyState) WaitResolving(name string) {
	status := d.GetStatus(name)
	<-status.done
}

// -------------------------------------- Data access methods --------------------------------------

func (d *DependencyState) SetFeatures(name FeatureGroupName, features []FeatureInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.features[name] = features
}

func (d *DependencyState) GetFeatures(name FeatureGroupName) []FeatureInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.features[name]
}

func (d *DependencyState) SetExceptions(name string, exceptions ExceptionInfo) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.exceptions[name] = exceptions
}

func (d *DependencyState) GetExceptions(name string) ExceptionInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.exceptions[name]
}

type FeatureInfo struct {
	EntityID string
	Value    int
}

type ExceptionInfo struct {
	Found bool
}

type FeatureGroupName string
