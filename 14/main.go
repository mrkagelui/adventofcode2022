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

	file, err := os.Open("14/input.txt")
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
	var lines []line

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ls, err := parseLines(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, ls...)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	g := withRocks(lines)

	log.Printf("success: %v\n", g.maxSand())
}

func part2(file *os.File) {
	var lines []line

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ls, err := parseLines(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, ls...)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	g := withRocksWithBottom(lines)

	log.Println("success!", g.maxSand())
}
