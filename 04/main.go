package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	switch r {
	case 1:
		part1()
	case 2:
		part2()
	}
}

func part1() {
	file, err := os.Open("04/input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	var num int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if oneFullyContainsAnother(parse(text)) {
			num++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success", num)
}

func part2() {
	file, err := os.Open("04/input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	var num int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if overlap(parse(text)) {
			num++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success", num)
}

func parse(s string) (int, int, int, int) {
	ss := strings.Split(s, `,`)
	if len(ss) != 2 {
		log.Fatalf("invalid string: %v", s)
	}
	s1s := strings.Split(ss[0], `-`)
	if len(s1s) != 2 {
		log.Fatalf("invalid string: %v, malformed s1: %v", s, ss[0])
	}
	min1, err := strconv.Atoi(s1s[0])
	if err != nil {
		log.Fatalf("invalid string: %v, malformed min1: %v", s, s1s[0])
	}
	max1, err := strconv.Atoi(s1s[1])
	if err != nil {
		log.Fatalf("invalid string: %v, malformed max1: %v", s, s1s[1])
	}

	s2s := strings.Split(ss[1], `-`)
	if len(s2s) != 2 {
		log.Fatalf("invalid string: %v, malformed s2: %v", s, ss[1])
	}
	min2, err := strconv.Atoi(s2s[0])
	if err != nil {
		log.Fatalf("invalid string: %v, malformed min2: %v", s, s2s[0])
	}
	max2, err := strconv.Atoi(s2s[1])
	if err != nil {
		log.Fatalf("invalid string: %v, malformed max2: %v", s, s2s[1])
	}

	return min1, max1, min2, max2
}

func oneFullyContainsAnother(min1, max1, min2, max2 int) bool {
	switch {
	case min1 < min2:
		return max2 <= max1
	case min1 == min2:
		return true
	default:
		return max2 >= max1
	}
}

func overlap(min1, max1, min2, max2 int) bool {
	return max1 >= min2 && min1 <= max2
}
