package cache

import (
	"MiniK8S/pkg/api/config"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient/listwatch"
	"errors"
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
	listerWatcher listwatch.ListerWatcher

	// MaxInternalErrorRetryDuration defines how long we should retry internal errors returned by watch.
	MaxInternalErrorRetryDuration time.Duration

	WorkQueue WorkQueue
}

func NewReflector(lw listwatch.ListerWatcher, expectedType apitypes.ApiObjectType, store Store, queue WorkQueue) *Reflector {
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
	err := r.ListAndWatch(stopCh, syncChan)
	if err != nil {
		panic(err)
		return
	}
}

func (r *Reflector) ListAndWatch(stopCh <-chan struct{}, syncChan chan bool) error {
	list, err := r.listerWatcher.List(config.ListOptions{
		Kind:            r.name,
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return err
	}
	err = r.HandleList(list)
	if err != nil {
		return err
	}
	syncChan <- true

	w, err := r.listerWatcher.Watch(config.ListOptions{
		Kind:            "",
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           true,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return err
	}

	err = r.HandleWatch(w, stopCh)
	w.Stop() // stop watch

	return nil
}

func (r *Reflector) HandleWatch(w watch.Interface, stopCh <-chan struct{}) error {
	for true {
		select {
		case <-stopCh:
			return errors.New("watch stopped")
		case event := <-w.ResultChan():
			switch event.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				r.PushEvent(event)

			case watch.Error:
				return errors.New("watch error")
			case watch.Bookmark:

				return errors.New("to be done")
			default:
				panic("should never get here")
				return errors.New("unknown watch event")
			}
		}
	}
	return nil
}

func (r *Reflector) HandleList(l config.ApiObjectList) error {
	list := l.GetItems()

	for _, obj := range list {
		key := obj.GetUID()
		r.store.Add(key.String(), obj)
	}
	return nil
}

func (r *Reflector) PushEvent(watchEvent watch.Event) {
	r.WorkQueue.Add(watchEvent)
}
