package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	file, err := os.Open("10/input.txt")
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
	ops := make([]*int, 0, 200)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ops = append(ops, parse(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Printf("success: %v", sumRegisterAtIntervals(ops, []int{20, 60, 100, 140, 180, 220}))
}

func part2(file *os.File) {
	ops := make([]*int, 0, 200)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ops = append(ops, parse(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success!")
	show(ops)
}

func parse(s string) *int {
	if strings.HasPrefix(s, "noop") {
		return nil
	}
	if len(s) < 6 {
		log.Fatal("invalid line", s)
	}
	n, err := strconv.Atoi(s[5:])
	if err != nil {
		log.Fatalf("not a num [%v]: %v", s, err)
	}
	return &n
}

func sumRegisterAtIntervals(ops []*int, intervals []int) int {
	var opIndex, interval, register, sum int
	register = 1
	for _, i := range intervals {
		for ; opIndex < len(ops); opIndex++ {
			cycles := 2
			if ops[opIndex] == nil {
				cycles = 1
			}
			if interval <= i && interval+cycles > i {
				sum += i * register
				break
			}

			interval += cycles
			if ops[opIndex] != nil {
				register += *ops[opIndex]
			}
		}
	}
	return sum
}

func show(ops []*int) {
	register := 1
	var cyclesSpentInOp, opIndex int
	for i := 0; i < 240; i++ {
		if i%40 == 0 {
			fmt.Println()
		}
		fmt.Printf("%c", char(i, register))

		if opIndex < len(ops) {
			cyclesSpentInOp++
			cycles := 1
			if ops[opIndex] != nil {
				cycles = 2
			}
			if cyclesSpentInOp == cycles {
				cyclesSpentInOp = 0
				if ops[opIndex] != nil {
					register += *ops[opIndex]
				}
				opIndex++
			}
		}
	}
}

func char(curr, register int) byte {
	curr = curr % 40
	if curr == register-1 || curr == register || curr == register+1 {
		return '#'
	}
	return '.'
}
