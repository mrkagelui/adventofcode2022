package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	var r int

	flag.IntVar(&r, "run", 1, "part to run")

	switch r {
	case 1:
		part1()
	case 2:
		part2()
	}
}

func part1() {
	file, err := os.Open("03/input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	var point int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		point += priority(text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success", point)
}

func part2() {
	file, err := os.Open("03/input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	var i, point int
	var s1, s2 string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		switch i {
		case 0:
			s1 = text
			i++
		case 1:
			s2 = text
			i++
		case 2:
			point += common(s1, s2, text)
			s1, s2, i = "", "", 0
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success", point)
}

func common(s1, s2, s3 string) int {
	return log2(hash(s1)&hash(s2)&hash(s3)) + 1
}

func priority(s string) int {
	return log2(hash(s[:len(s)/2])&hash(s[len(s)/2:])) + 1
}

func hash(s string) int64 {
	var r int64
	for _, c := range s {
		r |= code(c)
	}
	return r
}

func code(c rune) int64 {
	switch {
	case c >= 'A' && c <= 'Z':
		return 1 << (26 + (c - 'A'))
	default:
		return 1 << (c - 'a')
	}
}

func log2(input int64) int {
	var i int
	for ; input != 1; i++ {
		input = input >> 1
	}
	return i
}
