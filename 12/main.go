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

	file, err := os.Open("12/input.txt")
	//file, err := os.Open("12/test.txt")
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
	var m [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m = append(m, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	i, err := minimalSteps(m)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("success: %v", i)
}

func part2(file *os.File) {
	var m [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m = append(m, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	i, err := shortestElevatedPath(m)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success!", i)
}
