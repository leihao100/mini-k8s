package cache

import (
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient/listWatcher"
	"time"
)

type Reflector struct {
	// name identifies this reflector. By default it will be a file:line if possible.
	name string
	// The name of the type we expect to place in the store. The name
	// will be the stringification of expectedGVK if provided, and the
	// stringification of expectedType otherwise. It is for display
	// only, and should not be used for parsing or comparison.
	typeDescription string
	// An example object of the type we expect to place in the store.
	// Only the type needs to be right, except that when that is
	// `unstructured.Unstructured` the object's `"apiVersion"` and
	// `"kind"` must also be right.
	expectedType apitypes.ApiObjectType
	// The destination to sync up with the watch source
	store Store
	// listerWatcher is used to perform lists and watches.
	listerWatcher listWatcher.ListWatcher

	// MaxInternalErrorRetryDuration defines how long we should retry internal errors returned by watch.
	MaxInternalErrorRetryDuration time.Duration

	WorkQueue WorkQueue
}

func NewReflector(lw listWatcher.ListWatcher, expectedType apitypes.ApiObjectType, store Store, queue WorkQueue) *Reflector {
	return &Reflector{
		name:            string(expectedType) + "Reflector",
		typeDescription: string(expectedType),
		expectedType:    expectedType,
		listerWatcher:   lw,
		store:           store,
		WorkQueue:       queue,
	}
}

func (r *Reflector) Run(stopCh <-chan struct{}, syncChan chan bool) {

}

func (r *Reflector) ListAndWatch() {

}

func (r *Reflector) Watch() {

}

func (r *Reflector) HandleWatch() {

}

func (r *Reflector) List() {

}

func (r *Reflector) HandleList() {

}
