package class

import (
	"container/list"
)

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (queue *Queue) Enqueue(value interface{}) {
	queue.list.PushBack(value)
}

func (queue *Queue) Dequeue() (value interface{}, res bool) {
	if queue.Empty() {
		return value, false
	}
	value = queue.list.Front()
	queue.list.Remove(queue.list.Front())
	return value, true
}

func (queue *Queue) Head() (value interface{}, res bool) {
	if queue.Empty() {
		return value, false
	}
	value = queue.list.Front()
	return value, true
}

func (queue *Queue) Empty() bool {
	return queue.list.Len() == 0
}

func (queue *Queue) Size() int {
	return queue.list.Len()
}

func (queue *Queue) Clear() {
	for queue.Size() != 0 {
		queue.list.Remove(queue.list.Front())
	}
}
