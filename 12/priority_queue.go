package main

import "container/heap"

type priorityQueue struct {
	nodes  *nodeHeap
	handle map[int]*node // index of this does not change
}

func newPriorityQueue(points []point) priorityQueue {
	h := make(nodeHeap, len(points))
	handle := make(map[int]*node, len(points))
	for i, p := range points {
		n := node{
			p:        p,
			distance: nil,
			index:    i,
		}
		h[i] = &n
		handle[i] = &n
	}
	heap.Init(&h)
	return priorityQueue{
		nodes:  &h,
		handle: handle,
	}
}

// update modifies the priority corresponding to the index in the queue.
func (pq *priorityQueue) update(p point, distance int) {
	pointPtr := pq.handle[int(p)]
	pointPtr.distance = &distance
	heap.Fix(pq.nodes, pointPtr.index)
}

func (pq *priorityQueue) len() int {
	return pq.nodes.Len()
}

func (pq *priorityQueue) pop() any {
	n := heap.Pop(pq.nodes)
	nd := n.(*node)
	delete(pq.handle, int(nd.p))
	return n
}

func (pq *priorityQueue) contains(p point) bool {
	_, ok := pq.handle[int(p)]
	return ok
}

// ---- nodeHeap ----
type nodeHeap []*node

func (pq *nodeHeap) Len() int { return len(*pq) }

func (pq *nodeHeap) Less(i, j int) bool {
	if (*pq)[i].distance == nil {
		return false
	}
	if (*pq)[j].distance == nil {
		return (*pq)[i].distance != nil
	}
	return *(*pq)[i].distance < *(*pq)[j].distance
}

func (pq *nodeHeap) Swap(i, j int) {
	// the index of the node in heap changes at this step.
	// if node does not contain index, correspondence
	// of the index in heap (for fixing the heap) and
	// the index of the point cannot be established.
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *nodeHeap) Push(x any) {
	n := len(*pq)
	item := x.(*node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *nodeHeap) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type node struct {
	p        point
	distance *int
	index    int // see nodeHeap.Swap
}
