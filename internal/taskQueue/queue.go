package task_queue

import (
	"errors"

	"github.com/madchin/trader-bot/internal/gateway"
)

var ErrQueueIsEmpty = errors.New("queue is empty")

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

var Queue = queue[gateway.InteractionData]{}

func (q *queue[T]) Enqueue(data T) {
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

func (q *queue[T]) Dequeue() (T, error) {
	var data T
	if q.size == 0 {
		return data, ErrQueueIsEmpty
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

func (q *queue[T]) Size() int {
	return q.size
}
