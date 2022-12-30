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

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	file, err := os.Open("07/input.txt")
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
	const threshold = 100_000

	scanner := bufio.NewScanner(file)
	root := parse(scanner)

	root.calculateSize()
	total := root.sumSmallDirSize(threshold)

	log.Printf("success: %v\n", total)
}

func part2(file *os.File) {
	const total, needed = 70000000, 30000000
	scanner := bufio.NewScanner(file)
	root := parse(scanner)

	root.calculateSize()

	log.Printf("success: %v\n", root.sizeOfSmallestDirAbove(needed-(total-root.size)))
}

func parse(scanner *bufio.Scanner) *node {
	virtual := node{
		name: "root",
		size: 0,
		subs: []node{
			{
				name: "/",
			},
		},
	}
	curr := &virtual

	for scanner.Scan() {
		switch line := scanner.Text(); {
		case strings.HasPrefix(line, "$"): // a command
			switch command := line[2:]; {
			case command == "ls": // list content, can skip
				continue
			case strings.HasPrefix(command, "cd "): // change directory
				switch target := command[3:]; target {
				case "..": // go to parent
					curr = curr.parent
				default: // go to one sub
					i := search(curr.subs, target)
					if i == -1 {
						log.Fatalf("no such sub [%v] at [%v]", target, curr.name)
					}
					curr = &curr.subs[i]
				}
			}
		case strings.HasPrefix(line, "dir"): // content (dir)
			dirName := line[4:]
			if i := search(curr.subs, dirName); i != -1 {
				continue // already exists
			}
			curr.subs = append(curr.subs, node{
				name:   dirName,
				parent: curr,
			})
		default: // content (file)
			ss := strings.Split(line, " ")
			if len(ss) != 2 {
				log.Fatal("invalid file line", line)
			}
			sizeStr, name := ss[0], ss[1]
			if i := search(curr.subs, name); i != -1 {
				continue // already exists
			}
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				log.Fatalf("invalid size [%v] in line [%v]: %v", ss[0], line, err)
			}
			curr.subs = append(curr.subs, node{
				name:   name,
				size:   size,
				parent: curr,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	return &virtual.subs[0]
}

func search(nodes []node, name string) int {
	if len(nodes) == 0 {
		return -1
	}

	lo, hi := 0, len(nodes)
	for lo < hi {
		mid := (lo + hi) >> 1
		switch {
		case nodes[mid].name == name:
			return mid
		case nodes[mid].name < name:
			lo = mid + 1
		default:
			hi = mid
		}
	}
	return -1
}

// ---- tree ----

type node struct {
	name   string
	size   int
	subs   []node
	parent *node
}

func (n *node) isLeaf() bool {
	return n.subs == nil
}

func (n *node) calculateSize() {
	s := stack[*node]{}
	s.push(n)

	for {
		curr, err := s.peek()
		if err != nil {
			break
		}

		if !curr.isLeaf() && curr.size == 0 { // dir with unknown size
			for i := len(curr.subs) - 1; i >= 0; i-- {
				s.push(&curr.subs[i])
			}
			continue
		}
		if curr.parent != nil {
			curr.parent.size += curr.size
		}
		s.pop()
	}
}

func (n *node) sumSmallDirSize(threshold int) int {
	s := stack[*node]{}
	s.push(n)

	var result int
	for {
		curr, err := s.pop()
		if err != nil {
			break
		}
		if curr.size <= threshold {
			result += curr.size
		}
		for i := range curr.subs {
			if !curr.subs[i].isLeaf() {
				s.push(&curr.subs[i])
			}
		}
	}
	return result
}

func (n *node) sizeOfSmallestDirAbove(threshold int) int {
	s := stack[*node]{}
	s.push(n)

	var result int
	for {
		curr, err := s.pop()
		if err != nil {
			break
		}
		if result == 0 || curr.size >= threshold && curr.size < result {
			result = curr.size
		}
		for i := range curr.subs {
			if !curr.subs[i].isLeaf() {
				s.push(&curr.subs[i])
			}
		}
	}
	return result
}

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

func (s *stack[T]) isEmpty() bool {
	return len(*s) == 0
}
