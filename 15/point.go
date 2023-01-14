package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	row, col int
}

type pair struct {
	sensor, beacon point
}

func parse(s string) (pair, error) {
	// s is of format "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d"
	ss := strings.Split(s, ": closest beacon is at ")
	if len(ss) != 2 {
		return pair{}, fmt.Errorf("invalid str [%v]", s)
	}
	sensor, err := parsePoint(ss[0][10:])
	if err != nil {
		return pair{}, fmt.Errorf("parsing sensor in [%v]: %w", s, err)
	}
	beacon, err := parsePoint(ss[1])
	if err != nil {
		return pair{}, fmt.Errorf("parsing beacon in [%v]: %w", s, err)
	}
	return pair{sensor: sensor, beacon: beacon}, nil
}

func parsePoint(s string) (point, error) {
	// s is of format "x=%d, y=%d"
	i := strings.Index(s, ", y=")
	if i < 0 {
		return point{}, fmt.Errorf("invalid str [%v]", s)
	}
	xStr, yStr := s[2:i], s[i+4:]
	row, err := strconv.Atoi(yStr)
	if err != nil {
		return point{}, fmt.Errorf("parsing row in [%v]: %w", s, err)
	}
	col, err := strconv.Atoi(xStr)
	if err != nil {
		return point{}, fmt.Errorf("parsing col in [%v]: %w", s, err)
	}
	return point{row, col}, nil
}

func impossiblePoints(ps []pair, row int) int {
	seqs := make([]sequence, 0, len(ps))
	beaconsByRow := make(map[int]map[int]struct{}, len(ps))

	for _, p := range ps {
		if _, ok := beaconsByRow[p.beacon.row]; !ok {
			beaconsByRow[p.beacon.row] = make(map[int]struct{})
		}
		beaconsByRow[p.beacon.row][p.beacon.col] = struct{}{}
		seq, ok := intersections(p, row)
		if ok {
			seqs = append(seqs, seq)
		}
	}

	if len(seqs) == 0 {
		return 0
	}

	sort.Slice(seqs, func(i, j int) bool {
		return seqs[i].min < seqs[j].min
	})

	var r int
	curr := seqs[0]
	for i := 1; i < len(seqs); i++ {
		temp, ok := merge(curr, seqs[i])
		if !ok {
			r += curr.count()
		}
		curr = temp
	}
	r += curr.count()
	r -= len(beaconsByRow[row])

	return r
}

func intersections(p pair, row int) (sequence, bool) {
	radius := distance(p.sensor, p.beacon)
	centerCol, centerRow, upper, lower := p.sensor.col, p.sensor.row, p.sensor.row-radius, p.sensor.row+radius

	switch {
	case row < upper:
		return sequence{}, false
	case row >= upper && row <= centerRow:
		return sequence{centerCol - (row - upper), p.sensor.col + (row - upper)}, true
	case row > centerRow && row <= lower:
		return sequence{centerCol - (lower - row), centerCol + (lower - row)}, true
	default:
		return sequence{}, false
	}
}

func distance(a, b point) int {
	return diff(a.row, b.row) + diff(a.col, b.col)
}

func diff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

type sequence struct {
	min, max int
}

func (s sequence) count() int {
	return s.max - s.min + 1
}

// merge tries to combine two sequences a and b where a.min <= b.min, if not successful, return b
func merge(a, b sequence) (sequence, bool) {
	if b.min > a.max {
		return b, false
	}

	newMax := a.max
	if b.max > newMax {
		newMax = b.max
	}
	return sequence{
		min: a.min,
		max: newMax,
	}, true
}

func onlyBeaconFrequency(pairs []pair, maxRange int) int64 {
	var row, col int
	failedPoints := make(map[point]struct{})
	for _, pr := range pairs {
		p, ok := verifyPoints(shellOf(pr, maxRange), pairs, failedPoints)
		if ok {
			row = p.row
			col = p.col
			break
		}
	}
	return int64(col)*4_000_000 + int64(row)
}

func shellOf(p pair, max int) []point {
	radius := distance(p.sensor, p.beacon) + 1
	pts := make([]point, 0, radius*4-2)

	// lower right
	for i := 0; i < radius; i++ {
		col, row := p.sensor.col+(radius-i), p.sensor.row+i
		if col < 0 || col > max || row < 0 || row > max {
			continue
		}
		pts = append(pts, point{col: col, row: row})
	}

	// upper right
	for i := 0; i < radius; i++ {
		col, row := p.sensor.col+i, p.sensor.row-(radius-i)
		if col < 0 || col > max || row < 0 || row > max {
			continue
		}
		pts = append(pts, point{col: col, row: row})
	}

	// lower left
	for i := 0; i < radius; i++ {
		col, row := p.sensor.col-i, p.sensor.row+(radius-i)
		if col < 0 || col > max || row < 0 || row > max {
			continue
		}
		pts = append(pts, point{col: col, row: row})
	}

	// upper left
	for i := 0; i < radius; i++ {
		col, row := p.sensor.col-(radius-i), p.sensor.row-i
		if col < 0 || col > max || row < 0 || row > max {
			continue
		}
		pts = append(pts, point{col: col, row: row})
	}
	return pts
}

func verifyPoints(pts []point, pairs []pair, failedPoints map[point]struct{}) (point, bool) {
	for _, pt := range pts {
		if _, ok := failedPoints[pt]; ok {
			continue
		}
		if !isOutOfReach(pt, pairs) {
			failedPoints[pt] = struct{}{}
			continue
		}
		return pt, true
	}
	return point{}, false
}

func isOutOfReach(p point, pairs []pair) bool {
	for _, onePair := range pairs {
		if within(p, onePair) {
			return false
		}
	}
	return true
}

func within(p point, pr pair) bool {
	return distance(p, pr.sensor) <= distance(pr.sensor, pr.beacon)
}
