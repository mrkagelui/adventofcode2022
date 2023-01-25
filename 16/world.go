package main

import (
	"fmt"
	"math"

	"github.com/kagelui/adventofcode2022/priorityqueue"
)

type world map[string]valve

type node struct {
	prev         []string
	minutesSpent int
	released     int
}

func (n *node) visited(name string) bool {
	for _, s := range n.prev {
		if s == name {
			return true
		}
	}
	return false
}

func (w world) maxReleased(steps int) (int, error) {
	// find all working valves
	unopened := make(map[string]int)
	for s, v := range w {
		if v.flow > 0 {
			unopened[s] = v.flow
		}
	}

	_, released, err := w.maxReleasedSingle(steps, unopened)
	if err != nil {
		return 0, err
	}

	return released, nil
}

func (w world) maxReleasedSingle(steps int, unopened map[string]int) ([]string, int, error) {
	const starting = "AA"

	// cache stores all distances of other working valves from a particular valve
	cache := make(map[string]map[string]int, len(unopened))

	s := stack[node]{}
	s.push(node{
		prev:         []string{starting},
		minutesSpent: 0,
	})

	var valves []string
	var released int
	for !s.isEmpty() {
		// dfs to traverse all "nodes", which are the steps
		n, _ := s.pop()

		curr := n.prev[len(n.prev)-1]
		next, ok := cache[curr]
		if !ok {
			neighbors, err := w.distances(unopened, curr)
			if err != nil {
				return nil, 0, fmt.Errorf("dist: %w", err)
			}
			cache[curr] = neighbors
			next = neighbors
		}

		var hasNext bool
		for valveName, distance := range next {
			if t := steps - 1 - n.minutesSpent - distance; !n.visited(valveName) && t > 0 {
				s.push(node{
					prev:         attach(n.prev, valveName),
					minutesSpent: n.minutesSpent + distance + 1,
					released:     n.released + unopened[valveName]*t,
				})
				hasNext = true
			}
		}
		if !hasNext && n.released > released {
			valves, released = n.prev, n.released
		}
	}
	return valves, released, nil
}

func attach(ss []string, s string) []string {
	r := make([]string, len(ss)+1)
	for i, s2 := range ss {
		r[i] = s2
	}
	r[len(r)-1] = s
	return r
}

func (w world) distances(unopened map[string]int, start string) (map[string]int, error) {
	unreached := make(map[string]struct{}, len(unopened))
	for s := range unopened {
		unreached[s] = struct{}{}
	}
	result := make(map[string]int, len(unopened))

	all := make(map[string]int, len(w))
	for s := range w {
		all[s] = math.MaxInt
	}

	pq := priorityqueue.New(all)
	if err := pq.Update(start, 0); err != nil {
		return nil, fmt.Errorf("updating: %w", err)
	}

	for !pq.Empty() && len(unreached) != 0 {
		valveName, distance := pq.Pop()
		if distance == math.MaxInt {
			break // this means the rest is not reachable at all
		}
		if _, ok := unreached[valveName]; ok {
			result[valveName] = distance
			delete(unreached, valveName)
		}
		for _, v := range w[valveName].next {
			oldDistance, err := pq.Priority(v)
			if err != nil {
				continue
			}
			if d := distance + 1; d < oldDistance {
				if err := pq.Update(v, d); err != nil {
					return nil, fmt.Errorf("updating: %w", err)
				}
			}
		}
	}

	// clear any unreachable node from unopened
	for s := range unreached {
		delete(unopened, s)
	}

	return result, nil
}

func (w world) maxReleasedWithElephant(steps int) (int, error) {
	// find all working valves
	unopened := make(map[string]int)
	for s, v := range w {
		if v.flow > 0 {
			unopened[s] = v.flow
		}
	}

	meValves, meReleased, err := w.maxReleasedSingle(steps, unopened)
	if err != nil {
		return 0, err
	}
	for _, meValve := range meValves {
		delete(unopened, meValve)
	}
	_, eReleased, err := w.maxReleasedSingle(steps, unopened)
	if err != nil {
		return 0, err
	}

	return meReleased + eReleased, nil
}
