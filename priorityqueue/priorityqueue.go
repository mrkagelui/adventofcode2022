// Package priorityqueue provides a generic implementation of priority queue
// based on container/heap
package priorityqueue

import (
	"container/heap"
	"fmt"
)

type nodeHeap[K comparable, V ordered] []*node[K, V]

func (pq *nodeHeap[K, V]) Len() int { return len(*pq) }

func (pq *nodeHeap[K, V]) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority
}

func (pq *nodeHeap[K, V]) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *nodeHeap[K, V]) Push(x any) {
	item := x.(*node[K, V])
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *nodeHeap[K, V]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

type node[K comparable, V ordered] struct {
	value    K
	priority V
	index    int
}

type PriorityQueue[K comparable, V ordered] struct {
	nh     *nodeHeap[K, V]
	handle map[K]*node[K, V]
}

func New[K comparable, V ordered](m map[K]V) PriorityQueue[K, V] {
	var i int
	h := make(nodeHeap[K, V], len(m))
	handle := make(map[K]*node[K, V], len(m))
	for k, v := range m {
		n := &node[K, V]{
			value:    k,
			priority: v,
			index:    i,
		}
		h[i] = n
		i++

		handle[k] = n
	}
	heap.Init(&h)
	return PriorityQueue[K, V]{
		nh:     &h,
		handle: handle,
	}
}

func (pq *PriorityQueue[K, V]) Len() int {
	return pq.nh.Len()
}

func (pq *PriorityQueue[K, V]) Empty() bool {
	return pq.nh.Len() == 0
}

func (pq *PriorityQueue[K, V]) Update(k K, newPriority V) error {
	n, ok := pq.handle[k]
	if !ok {
		return fmt.Errorf("[%v] is not in the queue", k)
	}
	n.priority = newPriority
	heap.Fix(pq.nh, n.index)
	return nil
}

func (pq *PriorityQueue[K, V]) Pop() (K, V) {
	n := heap.Pop(pq.nh)
	nd := n.(*node[K, V])
	delete(pq.handle, nd.value)
	return nd.value, nd.priority
}

func (pq *PriorityQueue[K, V]) Push(k K, priority V) error {
	if _, ok := pq.handle[k]; ok {
		return fmt.Errorf("[%v] is already in queue", k)
	}
	n := &node[K, V]{value: k, priority: priority}
	heap.Push(pq.nh, n)
	pq.handle[k] = n
	return nil
}

func (pq *PriorityQueue[K, V]) Contains(k K) bool {
	_, ok := pq.handle[k]
	return ok
}

func (pq *PriorityQueue[K, V]) Priority(k K) (V, error) {
	n, ok := pq.handle[k]
	if !ok {
		return *new(V), fmt.Errorf("[%v] is not in the queue", k)
	}

	return n.priority, nil
}
