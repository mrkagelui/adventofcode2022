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

	file, err := os.Open("15/input.txt")
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
	var pairs []pair

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p, err := parse(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		pairs = append(pairs, p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success:", impossiblePoints(pairs, 2_000_000))
}

func part2(file *os.File) {
	var pairs []pair

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p, err := parse(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		pairs = append(pairs, p)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("starting...")

	log.Println("done!", onlyBeaconFrequency(pairs, 4_000_000))
}
