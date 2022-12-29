package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	file, err := os.Open("06/input.txt")
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
	const length = 4
	var input string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	i := firstIndexWithNDistinctRunes(input, length)
	log.Printf("success: %v, the runes are [%v]", i+length, input[i:i+length])
}

func part2(file *os.File) {
	const length = 14
	var input string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	i := firstIndexWithNDistinctRunes(input, length)
	log.Printf("success: %v, the runes are [%v]", i+length, input[i:i+length])
}

func firstIndexWithNDistinctRunes(s string, n int) int {
	var result int
	for i := 0; i < len(s)-n; {
		if incr := stepsToJump(s[i : i+n]); incr == 0 {
			result = i
			break
		} else {
			i += incr
		}
	}
	return result
}

func stepsToJump(s string) int {
	maxLastSeenIndex := -1
	lastSeenIndexes := []int{
		-1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1,
	}

	for i, c := range s {
		index := indexOf(c)
		lsi := lastSeenIndexes[index]
		if lsi == -1 {
			lastSeenIndexes[index] = i
			continue
		}
		if lsi > maxLastSeenIndex {
			maxLastSeenIndex = lsi
		}
	}
	return maxLastSeenIndex + 1
}

func indexOf(c rune) int {
	return int(c - 'a')
}
