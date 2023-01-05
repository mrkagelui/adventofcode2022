package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	file, err := os.Open("09/input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	switch r {
	case 1:
		part1(file)
	case 2:
		part2(file)
	}
}

func part1(file *os.File) {
	points := make([]point, 2)
	allTails := make(map[point]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m, t := parse(scanner.Text())
		walk(points, m, t, allTails)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Printf("success: %v", len(allTails))
}

func part2(file *os.File) {
	points := make([]point, 10)
	allTails := make(map[point]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m, t := parse(scanner.Text())
		walk(points, m, t, allTails)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Printf("success: %v", len(allTails))
}

func up(s point) point {
	return point{s.x, s.y - 1}
}

func down(s point) point {
	return point{s.x, s.y + 1}
}

func left(s point) point {
	return point{s.x - 1, s.y}
}

func right(s point) point {
	return point{s.x + 1, s.y}
}

func upperLeft(p point) point {
	return point{p.x - 1, p.y - 1}
}

func upperRight(p point) point {
	return point{p.x + 1, p.y - 1}
}

func lowerLeft(p point) point {
	return point{p.x - 1, p.y + 1}
}

func lowerRight(p point) point {
	return point{p.x + 1, p.y + 1}
}

type move func(point) point

func parse(s string) (move, int) {
	ss := strings.Split(s, " ")
	if len(ss) != 2 {
		log.Fatalf("invalid line [%v]", s)
	}

	times, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Fatalf("invalid times [%v]: %v", s, err)
	}

	switch ss[0] {
	case "U":
		return up, times
	case "D":
		return down, times
	case "L":
		return left, times
	case "R":
		return right, times
	default:
		log.Fatalf("invalid move [%v]: %v", s, err)
		return nil, 0
	}
}

func walk(points []point, m move, times int, records map[point]struct{}) {
	if len(points) == 0 {
		return
	}
	for i := 0; i < times; i++ {
		points[0] = m(points[0])
		for j := 1; j < len(points); j++ {
			var moved bool
			points[j], moved = follow(points[j-1], points[j])
			if !moved {
				break
			}
		}
		records[points[len(points)-1]] = struct{}{}
	}
}

func follow(head, tail point) (point, bool) {
	horizontalDistance, verticalDistance := head.x-tail.x, head.y-tail.y
	if abs(horizontalDistance) <= 1 && abs(verticalDistance) <= 1 {
		return tail, false
	}

	switch {
	case horizontalDistance < 0:
		switch {
		case verticalDistance < 0:
			return upperLeft(tail), true
		case verticalDistance == 0:
			return left(tail), true
		default:
			return lowerLeft(tail), true
		}
	case horizontalDistance == 0:
		if verticalDistance < 0 {
			return up(tail), true
		}
		return down(tail), true
	default:
		switch {
		case verticalDistance < 0:
			return upperRight(tail), true
		case verticalDistance == 0:
			return right(tail), true
		default:
			return lowerRight(tail), true
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
