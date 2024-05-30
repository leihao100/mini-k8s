package cache

import (
	"sync"
)

type WorkQueue struct {
	queue         []interface{}
	processingSet map[interface{}]struct{}
	mutex         sync.Mutex
	cond          *sync.Cond
	shutdown      bool
}

func NewWorkQueue() *WorkQueue {
	wq := &WorkQueue{
		queue:         make([]interface{}, 0),
		processingSet: make(map[interface{}]struct{}),
		shutdown:      false,
	}
	wq.cond = sync.NewCond(&wq.mutex)
	return wq
}

func (wq *WorkQueue) Add(item interface{}) {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	if _, exists := wq.processingSet[item]; exists {
		return
	}

	wq.queue = append(wq.queue, item)
	wq.cond.Signal()
}

func (wq *WorkQueue) Get() (interface{}, bool) {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	for len(wq.queue) == 0 && !wq.shutdown {
		wq.cond.Wait()
	}

	if wq.shutdown {
		return nil, true
	}

	item := wq.queue[0]
	wq.queue = wq.queue[1:]
	wq.processingSet[item] = struct{}{}
	return item, false
}

func (wq *WorkQueue) Done(item interface{}) {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	delete(wq.processingSet, item)
	if len(wq.queue) == 0 && wq.shutdown {
		wq.cond.Broadcast()
	}
}

func (wq *WorkQueue) ShutDown() {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	wq.shutdown = true
	wq.cond.Broadcast()
}

func (wq *WorkQueue) ShutDownWithDrain() {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	wq.shutdown = true
	for len(wq.queue) > 0 {
		wq.cond.Wait()
	}
	wq.cond.Broadcast()
}

func (wq *WorkQueue) Len() int {
	wq.mutex.Lock()
	defer wq.mutex.Unlock()

	return len(wq.queue)
}
