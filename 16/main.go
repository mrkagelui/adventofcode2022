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

	file, err := os.Open("16/input.txt")
	//file, err := os.Open("16/test.txt")
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
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	w, err := parse(lines)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting")
	r, err := w.maxReleased(30)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success:", r)
}

func part2(file *os.File) {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	w, err := parse(lines)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting")
	r, err := w.maxReleasedWithElephant(26)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success:", r)
}
