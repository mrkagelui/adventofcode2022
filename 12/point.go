package main

type point int

func canGo(s, e point, m matrix) bool {
	start := m.elevationAt(s)
	end := m.elevationAt(e)

	return end <= start+1
}

func canReverse(s, e point, m matrix) bool {
	start := m.elevationAt(s)
	end := m.elevationAt(e)

	return end >= start-1
}
