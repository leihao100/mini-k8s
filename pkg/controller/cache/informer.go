package cache

import (
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient/listWatcher"
)

type Informer struct {
	ty        types.ApiObjectType
	reflector *Reflector
	handlers  []EventHandler

	// store is shared with reflector, which
	// stores objects of objType get from ApiServer
	store Store

	// transportQueue is used to get notification from
	// reflector about new events happening
	queue WorkQueue
}

func NewInformer(ty types.ApiObjectType, store Store, queue WorkQueue, lw listWatcher.ListWatch, h EventHandler) *Informer {
	return &Informer{
		ty:        ty,
		queue:     queue,
		store:     store,
		reflector: NewReflector(lw, ty, store, queue),
		handlers:  []EventHandler{h},
	}
}

func (i *Informer) Run(stopCh <-chan struct{}) {

}

func (i *Informer) Get(key string) (item interface{}, exists bool, err error) {
	return i.store.Get(key)
}

func (i *Informer) List() []interface{} {
	return i.store.List()
}

type EventHandler interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
}

// EventHandlerFuncs is an implementation of EventHandler
type EventHandlerFuncs struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}

// OnAdd calls AddFunc if it's not nil.
func (r EventHandlerFuncs) OnAdd(obj interface{}) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r EventHandlerFuncs) OnUpdate(oldObj, newObj interface{}) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r EventHandlerFuncs) OnDelete(obj interface{}) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
}
