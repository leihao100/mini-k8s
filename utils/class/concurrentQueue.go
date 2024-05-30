package class

import (
	"container/list"
	"sync"
)

// ConcurrentQueue 是一个并发安全的队列类型
type ConcurrentQueue struct {
	list *list.List
	mu   sync.Mutex
}

// NewConcurrentQueue 创建一个新的并发安全队列
func NewConcurrentQueue() *ConcurrentQueue {
	return &ConcurrentQueue{
		list: list.New(),
	}
}

// Enqueue 将值添加到队列尾部
func (queue *ConcurrentQueue) Enqueue(value interface{}) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.list.PushBack(value)
}

// Dequeue 从队列头部移除并返回一个值
func (queue *ConcurrentQueue) Dequeue() (value interface{}, ok bool) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.list.Len() == 0 {
		return nil, false
	}
	elem := queue.list.Front()
	queue.list.Remove(elem)
	return elem.Value, true
}

// Front 返回队列头部的值但不移除它
func (queue *ConcurrentQueue) Front() (value interface{}, ok bool) {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.list.Len() == 0 {
		return nil, false
	}
	elem := queue.list.Front()
	return elem.Value, true
}

// Empty 返回队列是否为空
func (queue *ConcurrentQueue) Empty() bool {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	return queue.list.Len() == 0
}

// Size 返回队列的长度
func (queue *ConcurrentQueue) Size() int {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	return queue.list.Len()
}

// Clear 清空队列
func (queue *ConcurrentQueue) Clear() {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	queue.list.Init()
}
