package main

import (
	"fmt"
	"strconv"
	"strings"
)

type point struct {
	row, col int
}
type line struct {
	start, end point
}

func parseLines(s string) ([]line, error) {
	var res []line

	ptStrs := strings.Split(s, " -> ")
	if len(ptStrs) < 2 {
		return nil, fmt.Errorf("too few points in [%v]", s)
	}

	start, err := parsePoint(ptStrs[0])
	if err != nil {
		return nil, fmt.Errorf("parsing point: %v", err)
	}

	for i := 1; i < len(ptStrs); i++ {
		p, err := parsePoint(ptStrs[i])
		if err != nil {
			return nil, fmt.Errorf("parsing point: %v", err)
		}
		res = append(res, line{start: start, end: p})
		start = p
	}

	return res, nil
}

func parsePoint(s string) (point, error) {
	ss := strings.Split(s, ",")
	if len(ss) != 2 {
		return point{}, fmt.Errorf("invalid point [%v]", ss)
	}
	col, err := strconv.Atoi(ss[0])
	if err != nil {
		return point{}, fmt.Errorf("parsing col in [%v]: %v", ss, err)
	}
	row, err := strconv.Atoi(ss[1])
	if err != nil {
		return point{}, fmt.Errorf("parsing row in [%v]: %v", ss, err)
	}
	return point{
		row: row,
		col: col,
	}, nil
}
