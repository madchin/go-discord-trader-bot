package scheduler

import (
	"errors"
)

var errQueueIsEmpty = errors.New("queue is empty")

type node[T any] struct {
	data T
	prev *node[T]
	next *node[T]
}

type queue[T any] struct {
	first *node[T]
	last  *node[T]
	size  int
}

func (q *queue[T]) enqueue(data T) {
	node := &node[T]{data: data}
	if q.size == 0 {
		q.first = node
		q.last = q.first
	} else {
		q.last.next = node
		node.prev = q.last
		q.last = node
	}
	q.size++
}

func (q *queue[T]) dequeue() (T, error) {
	var data T
	if q.size == 0 {
		return data, errQueueIsEmpty
	}
	if q.size == 1 {
		data = q.first.data
		q.first = nil
		q.last = nil
	} else {
		data = q.first.data
		q.first.next.prev = nil
		q.first = q.first.next
	}
	q.size--
	return data, nil
}
