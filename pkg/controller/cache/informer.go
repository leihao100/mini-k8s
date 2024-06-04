package cache

import (
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	"fmt"
	"time"
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
	queue *WorkQueue
}

func NewInformer(ty types.ApiObjectType, store Store, queue *WorkQueue, lw listwatch.ListerWatcher, h EventHandler) *Informer {
	return &Informer{
		ty:        ty,
		queue:     queue,
		store:     store,
		reflector: NewReflector(lw, ty, store, queue),
		handlers:  []EventHandler{h},
	}
}

func NewDefaultInformerAndCli(ty types.ApiObjectType) (*apiClient.Client, *Informer) {
	cli := apiClient.NewRESTClient(ty)
	lw := listwatch.NewListWatchFromClient(cli)
	store := NewSimpleStore()
	queue := NewWorkQueue()
	ref := NewReflector(lw, ty, store, queue)
	return cli, &Informer{
		ty:        ty,
		reflector: ref,
		handlers:  []EventHandler{},
		store:     store,
		queue:     queue,
	}
}

func (i *Informer) Run(stopCh <-chan struct{}) {
	fmt.Println("[Informer ", i.reflector.expectedType, "]Informer starting...")
	syncChan := make(chan bool)

	go func() {
		i.reflector.Run(stopCh, syncChan)
	}()
	//waiting for list
	<-syncChan

	go func() {
		for {
			select {
			case <-stopCh:
				return

			default:
				if i.queue.Len() == 0 {
					time.Sleep(1 * time.Second)
					continue
				}
				obj, shutdown := i.queue.Get()
				if shutdown {
					continue
				}
				event, ok := obj.(watch.Event)
				if !ok {
					panic("informer translate object to watch.event failed")
				}
				fmt.Println("[informer]", i.reflector.expectedType, " translate object to watch.event]")
				switch event.Type {
				case watch.Added:
					i.store.Update(event.Object.GetUID().String(), event.Object)
					for _, h := range i.handlers {
						h.OnAdd(event.Object)
					}
				case watch.Modified:
					old, exist, _ := i.store.Get(event.Object.GetUID().String())
					i.store.Update(event.Object.GetUID().String(), event.Object)
					if !exist {
						for _, h := range i.handlers {
							h.OnAdd(event.Object)
						}
					} else {
						for _, h := range i.handlers {
							h.OnUpdate(old, event.Object)
						}
					}

				case watch.Deleted:
					fmt.Println("[informer]", i.reflector.expectedType, " translate delete into to watch.event")
					obj, exist, _ := i.store.Get(event.Object.GetUID().String())
					fmt.Println("[informer]", i.reflector.expectedType, " delete object's existence is", exist)
					if exist {
						err := i.store.Delete(event.Object.GetUID().String())
						if err != nil {
							return
						}
						for _, handler := range i.handlers {
							handler.OnDelete(obj)
						}
					}
				}
			}
		}
	}()

}

func (i *Informer) Get(key string) (item interface{}, exists bool, err error) {
	return i.store.Get(key)
}

func (i *Informer) List() []interface{} {
	return i.store.List()
}

func (i *Informer) AddEventHandler(h EventHandler) {
	i.handlers = append(i.handlers, h)
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
