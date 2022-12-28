package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	var first, second, third, current int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			first, second, third = rank(first, second, third, current)
			current = 0
			continue
		}
		i, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal(text, "not a number", err)
		}
		current += i
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success", first, second, third, "total", first+second+third)
}

func rank(first, second, third, n int) (int, int, int) {
	if n > first {
		return n, first, second
	}
	if n > second {
		return first, n, second
	}
	if n > third {
		return first, second, n
	}
	return first, second, third
}
