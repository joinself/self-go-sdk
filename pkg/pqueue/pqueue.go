// Copyright 2020 Self Group Ltd. All Rights Reserved.

package pqueue

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// Queue a priority queue implementation
type Queue struct {
	partitions []list
	items      int64
	maxLength  int64
	cond       *sync.Cond
}

// New creates a new prioriy queue
func New(priorities, maxLength int) *Queue {
	return &Queue{
		partitions: make([]list, priorities),
		maxLength:  int64(maxLength),
		cond:       sync.NewCond(&sync.Mutex{}),
	}
}

// Push an item to the queue
func (q *Queue) Push(priority int, value interface{}) {
	// fmt.Println(atomic.LoadInt64(&q.items))
	for atomic.LoadInt64(&q.items) > q.maxLength {
		runtime.Gosched()
	}

	q.partitions[priority].push(value)
	q.cond.L.Lock()
	q.cond.Signal()
	q.cond.L.Unlock()

	atomic.AddInt64(&q.items, 1)
}

// Pop an item from the queue
func (q *Queue) Pop() interface{} {
	_, v := q.pop()
	return v
}

// Pop an item from the queue with its priority
func (q *Queue) PopWithPrioriry() (int, interface{}) {
	return q.pop()
}

func (q *Queue) pop() (int, interface{}) {
	for i := 0; i < len(q.partitions); i++ {
		if q.partitions[i].empty() {
			continue
		}

		v := q.partitions[i].pop()
		if v != nil {
			atomic.AddInt64(&q.items, -1)
			return i, v
		}
	}

	q.cond.L.Lock()
	q.cond.Wait()
	q.cond.L.Unlock()

	return q.pop()
}

// Flush clears a partitions results
func (q *Queue) Flush(priority int) {
	q.partitions[priority].flush()
}
