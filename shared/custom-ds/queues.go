package customds

import (
	"sync"
)

// When a goroutine tries to push to a full queue or pull from an empty queue, it will get blocked.
type BlockingQueue[T any] struct {
	m        sync.Mutex
	c        sync.Cond
	data     []T
	capacity int
}

func NewBlockingBlockingQueue[T any](capacity int) *BlockingQueue[T] {
	BlockingQueue := new(BlockingQueue[T])
	BlockingQueue.c = sync.Cond{L: &BlockingQueue.m}
	BlockingQueue.capacity = capacity

	return BlockingQueue
}

func (q *BlockingQueue[T]) isFull() bool {
	return len(q.data) == q.capacity
}

func (q *BlockingQueue[T]) isEmpty() bool {
	return len(q.data) == 0
}

func (q *BlockingQueue[T]) Put(input T) {
	q.c.L.Lock()

	defer q.c.L.Unlock()

	if q.isFull() {
		q.c.Wait()
	}

	q.data = append(q.data, input)

	q.c.Signal()
}

func (q *BlockingQueue[T]) Take() T {
	q.c.L.Lock()
	defer q.c.L.Unlock()

	if q.isEmpty() {
		q.c.Wait()
	}

	result := q.data[0]

	q.data = q.data[1:len(q.data)]
	q.c.Signal()
	return result
}

/*
// Implementation using channels
type BlockingQueue struct {
	channel chan interface{}
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	q := &BlockingQueue{make(chan interface{}, capacity)}

	return q
}

func (q *BlockingQueue ) Put(item interface{}) {
	q.channel <- item
}

func (q *BlockingQueue ) Take() interface{} {
	return <-q.channel
}
*/
