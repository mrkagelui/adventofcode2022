package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// ---- stack ----
type stack[T any] []T

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s *stack[T]) pop() (T, error) {
	if len(*s) == 0 {
		return *new(T), errors.New("empty stack")
	}
	v := []T(*s)[len(*s)-1]
	*s = []T(*s)[:len(*s)-1]
	return v, nil
}

func (s *stack[T]) peek() (T, error) {
	if len(*s) == 0 {
		return *new(T), errors.New("empty stack")
	}
	return []T(*s)[len(*s)-1], nil
}

// ---- end of stack ----

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	file, err := os.Open("05/input.txt")
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
	var dLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}
		dLines = append(dLines, text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan drawing", err)
	}

	stacks := parseDrawing(dLines)

	for scanner.Scan() {
		text := scanner.Text()
		times, from, to := parseMove(text)
		move(stacks, times, from-1, to-1)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan moves", err)
	}

	var bytes []byte
	for i, s := range stacks {
		top, err := s.peek()
		if err != nil {
			log.Fatalf("stack %v [%v]: %v", i, s, err)
		}
		bytes = append(bytes, top)
	}

	log.Println("success", string(bytes))
}

func part2(file *os.File) {
	var dLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}
		dLines = append(dLines, text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan drawing", err)
	}

	stacks := parseDrawing(dLines)

	for scanner.Scan() {
		text := scanner.Text()
		times, from, to := parseMove(text)
		moveMultiple(stacks, times, from-1, to-1)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan moves", err)
	}

	var bytes []byte
	for i, s := range stacks {
		top, err := s.peek()
		if err != nil {
			log.Fatalf("stack %v [%v]: %v", i, s, err)
		}
		bytes = append(bytes, top)
	}

	log.Println("success", string(bytes))
}

func parseDrawing(lines []string) []stack[byte] {
	numStacks := (len(lines[0]) + 1) / 4

	matrix := make([][]byte, numStacks)

	for i := len(lines) - 2; i >= 0; i-- {
		for j := 0; j < numStacks; j++ {
			if r := lines[i][4*j+1]; r != ' ' {
				matrix[j] = append(matrix[j], r)
			}
		}
	}

	stacks := make([]stack[byte], len(matrix))
	for i, bytes := range matrix {
		stacks[i] = bytes
	}

	return stacks
}

func parseMove(s string) (int, int, int) {
	ss := strings.Split(s, ` `)
	if len(ss) != 6 {
		log.Fatalf("invalid move: %v", s)
	}
	times, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Fatalf("[%v] contains invalid times: %v", s, err)
	}
	from, err := strconv.Atoi(ss[3])
	if err != nil {
		log.Fatalf("[%v] contains invalid from: %v", s, err)
	}
	to, err := strconv.Atoi(ss[5])
	if err != nil {
		log.Fatalf("[%v] contains invalid to: %v", s, err)
	}
	return times, from, to
}

func move[T any](stacks []stack[T], n, from, to int) {
	for i := 0; i < n; i++ {
		v, err := stacks[from].pop()
		if err != nil {
			log.Fatalf("moving %v times from %v to %v for stack %v at %v", n, from, to, stacks, i)
		}
		stacks[to].push(v)
	}
}

func moveMultiple[T any](stacks []stack[T], n, from, to int) {
	temp := stack[T]{}
	for i := 0; i < n; i++ {
		v, err := stacks[from].pop()
		if err != nil {
			log.Fatalf("moving %v times from %v to temp for stack %v at %v", n, from, stacks, i)
		}
		temp.push(v)
	}
	for i := 0; i < n; i++ {
		v, err := temp.pop()
		if err != nil {
			log.Fatalf("moving %v times from temp to %v for stack %v at %v", n, to, stacks, i)
		}
		stacks[to].push(v)
	}
}
