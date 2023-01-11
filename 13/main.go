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

	file, err := os.Open("13/input.txt")
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
	var sum int
	var n node
	var firstNodeRead bool
	i := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if !firstNodeRead {
			var err error
			n, err = parseNode(line)
			if err != nil {
				log.Fatal(err)
			}
			firstNodeRead = true
			continue
		}

		m, err := parseNode(line)
		if err != nil {
			log.Fatal(err)
		}
		smallerPtr := smallerThan(n, m)
		if smallerPtr == nil {
			log.Fatal("same node in pair", i)
		}
		if *smallerPtr {
			sum += i
		}

		firstNodeRead = false
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Printf("success: %v\n", sum)
}

func part2(file *os.File) {
	a := node{subs: []node{{subs: []node{{value: 2}}}}}
	b := node{subs: []node{{subs: []node{{value: 6}}}}}
	var smallerThanA, smallerThanB int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		n, err := parseNode(line)
		if err != nil {
			log.Fatal(err)
		}

		s := smallerThan(n, a)
		if s == nil {
			log.Fatal("same node: ", line)
		}
		if *s {
			smallerThanA++
		}
		s = smallerThan(n, b)
		if s == nil {
			log.Fatal("same node: ", line)
		}
		if *s {
			smallerThanB++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	result := (smallerThanA + 1) * (smallerThanB + 2)
	if smallerThanA > smallerThanB {
		result = (smallerThanA + 2) * (smallerThanB + 1)
	}

	log.Println("success!", result)
}
