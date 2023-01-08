package main

import "errors"

type matrix [][]byte

func (m matrix) leftOf(p point) (point, bool) {
	if len(m) == 0 {
		return 0, false
	}
	width := len(m[0])
	if int(p)%width <= 0 {
		return 0, false
	}
	return p - 1, true
}

func (m matrix) rightOf(p point) (point, bool) {
	if len(m) == 0 {
		return 0, false
	}
	width := len(m[0])
	if int(p)%width >= width-1 {
		return 0, false
	}
	return p + 1, true
}

func (m matrix) upOf(p point) (point, bool) {
	if len(m) == 0 {
		return 0, false
	}
	width := len(m[0])

	if int(p)/width <= 0 {
		return 0, false
	}
	return point(int(p) - width), true
}

func (m matrix) downOf(p point) (point, bool) {
	if len(m) == 0 {
		return 0, false
	}
	width := len(m[0])
	height := len(m)

	if int(p)/width >= height-1 {
		return 0, false
	}
	return point(int(p) + width), true
}

func (m matrix) elevationAt(p point) byte {
	width := len(m[0])
	row := int(p) / width
	col := int(p) % width

	e := m[row][col]

	switch e {
	case 'S':
		e = 'a'
	case 'E':
		e = 'z'
	}

	return e
}

// minimalSteps gives the length of the shortest path from starting point S to ending point E
func minimalSteps(m matrix) (int, error) {
	var s, e point
	height, width := len(m), len(m[0])

	pts := make([]point, width*height)
	for i, row := range m {
		for j, ch := range row {
			index := i*width + j
			p := point(index)
			switch ch {
			case 'S':
				s = p
			case 'E':
				e = p
			}

			pts[index] = p
		}
	}
	q := newPriorityQueue(pts)
	q.update(s, 0)

	for q.len() > 0 {
		n := q.pop().(*node)
		if n.distance == nil {
			return 0, errors.New("no more connected node")
		}
		if n.p == e {
			return *n.distance, nil
		}

		l, ok := m.leftOf(n.p)
		if ok && q.contains(l) && canGo(n.p, l, m) {
			q.update(l, *n.distance+1)
		}

		r, ok := m.rightOf(n.p)
		if ok && q.contains(r) && canGo(n.p, r, m) {
			q.update(r, *n.distance+1)
		}

		u, ok := m.upOf(n.p)
		if ok && q.contains(u) && canGo(n.p, u, m) {
			q.update(u, *n.distance+1)
		}

		d, ok := m.downOf(n.p)
		if ok && q.contains(d) && canGo(n.p, d, m) {
			q.update(d, *n.distance+1)
		}
	}
	return 0, errors.New("no more nodes")
}

// shortestElevatedPath gives the length of the shortest elevated path.
// an elevated path is from any point of elevation 'a' to the ending point E.
// this is achieved by starting from the point E and stopping when any point at 'a' is reached
func shortestElevatedPath(m matrix) (int, error) {
	var s point
	height, width := len(m), len(m[0])

	pts := make([]point, width*height)
	for i, row := range m {
		for j, ch := range row {
			index := i*width + j

			p := point(index)
			if ch == 'E' {
				s = p
			}
			pts[index] = p
		}
	}
	q := newPriorityQueue(pts)
	q.update(s, 0)

	for q.len() > 0 {
		n := q.pop().(*node)
		if n.distance == nil {
			return 0, errors.New("no more connected node")
		}
		if m.elevationAt(n.p) == 'a' {
			return *n.distance, nil
		}

		l, ok := m.leftOf(n.p)
		if ok && q.contains(l) && canReverse(n.p, l, m) {
			q.update(l, *n.distance+1)
		}

		r, ok := m.rightOf(n.p)
		if ok && q.contains(r) && canReverse(n.p, r, m) {
			q.update(r, *n.distance+1)
		}

		u, ok := m.upOf(n.p)
		if ok && q.contains(u) && canReverse(n.p, u, m) {
			q.update(u, *n.distance+1)
		}

		d, ok := m.downOf(n.p)
		if ok && q.contains(d) && canReverse(n.p, d, m) {
			q.update(d, *n.distance+1)
		}
	}
	return 0, errors.New("no more nodes")
}
