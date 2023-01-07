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

type item struct {
	level int
}

func (i *item) update(newLevel int) {
	i.level = newLevel
}

type delivery struct {
	i  item
	to int
}

type monkey struct {
	items    []item
	op       func(*item)
	checkNum int
	yesTo    int
	noTo     int
	num      int
	version  int
	bound    int
}

func (m *monkey) receive(i item) {
	m.items = append(m.items, i)
}

func (m *monkey) play() []delivery {
	deliveries := make([]delivery, len(m.items))
	for i := range m.items {
		m.num++

		m.op(&m.items[i])
		switch m.version {
		case 1:
			m.items[i].level /= 3
		default:
			m.items[i].level %= m.bound
		}

		to := m.noTo
		if m.items[i].level%m.checkNum == 0 {
			to = m.yesTo
		}

		deliveries[i] = delivery{
			i:  m.items[i],
			to: to,
		}
	}
	m.items = nil

	return deliveries
}

func deliver(all []monkey, deliveries []delivery) error {
	for _, d := range deliveries {
		if d.to < 0 || d.to >= len(all) {
			return fmt.Errorf("no such monkey: %v", d.to)
		}
		all[d.to].receive(d.i)
	}
	return nil
}

func main() {
	var r int

	flag.IntVar(&r, "run", 2, "part to run")

	file, err := os.Open("11/input.txt")
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
	const (
		rounds    = 20
		numConfig = 6
		version   = 1
	)
	all := make([]monkey, 0)

	var config []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line == "" {
			continue
		} else {
			config = append(config, line)
		}

		if len(config) == numConfig {
			m := parseMonkey(config, version)
			all = append(all, m)
			config = nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	for i := 0; i < rounds; i++ {
		for j := range all {
			if err := deliver(all, all[j].play()); err != nil {
				log.Fatal(err)
			}
		}
	}

	var top, second int
	for _, m := range all {
		if m.num > top {
			second = top
			top = m.num
			continue
		}
		if m.num > second {
			second = m.num
		}
	}

	log.Printf("success: %v", top*second)
}

func part2(file *os.File) {
	const (
		rounds    = 10000
		numConfig = 6
		version   = 2
	)
	all := make([]monkey, 0)

	var config []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line == "" {
			continue
		} else {
			config = append(config, line)
		}

		if len(config) == numConfig {
			m := parseMonkey(config, version)
			all = append(all, m)
			config = nil
		}
	}

	bound := 1
	for _, m := range all {
		bound *= m.checkNum
	}
	for i := range all {
		all[i].bound = bound
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	for i := 0; i < rounds; i++ {
		for j := range all {
			if err := deliver(all, all[j].play()); err != nil {
				log.Fatal(err)
			}
		}
	}

	var top, second int
	for _, m := range all {
		if m.num > top {
			second = top
			top = m.num
			continue
		}
		if m.num > second {
			second = m.num
		}
	}

	log.Printf("success: %v", top*second)
}

func parseMonkey(contents []string, version int) monkey {
	var items []item
	var op func(*item)
	var checkNum, yesTo, noTo int

	for _, content := range contents {
		switch content = strings.Trim(content, " "); {
		case strings.HasPrefix(content, "Monkey"), content == "":
			continue
		case strings.HasPrefix(content, "Starting items: "):
			items = parseItems(strings.Split(content[len("Starting items: "):], ", "))
		case strings.HasPrefix(content, "Operation: "):
			ss := strings.Split(content, " ")
			op = parseOp(ss[len(ss)-2:])
		case strings.HasPrefix(content, "Test: divisible by "):
			n, err := strconv.Atoi(content[len("Test: divisible by "):])
			if err != nil {
				log.Fatalf("err check [%v]: %v", content, err)
			}
			checkNum = n
		case strings.HasPrefix(content, "If true: throw to monkey "):
			n, err := strconv.Atoi(content[len("If true: throw to monkey "):])
			if err != nil {
				log.Fatalf("err YesTo [%v]: %v", content, err)
			}
			yesTo = n
		case strings.HasPrefix(content, "If false: throw to monkey "):
			n, err := strconv.Atoi(content[len("If false: throw to monkey "):])
			if err != nil {
				log.Fatalf("err NoTo [%v]: %v", content, err)
			}
			noTo = n
		default:
			log.Printf("ERR config [%v]\n", content)
		}
	}
	return monkey{
		items:    items,
		op:       op,
		checkNum: checkNum,
		yesTo:    yesTo,
		noTo:     noTo,
		version:  version,
	}
}

func parseItems(levels []string) []item {
	items := make([]item, len(levels))
	for i, levelStr := range levels {
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			log.Fatalf("not a number [%v]: %v", levelStr, err)
		}
		items[i] = item{level: level}
	}
	return items
}

func parseOp(ss []string) func(*item) {
	if len(ss) != 2 {
		log.Fatalf("len not 2: [%v]", ss)
	}

	if ss[1] == "old" {
		switch ss[0] {
		case "+":
			return func(i *item) {
				i.level += i.level
			}
		case "-":
			return func(i *item) {
				i.level = 0
			}
		case "*":
			return func(i *item) {
				i.level *= i.level
			}
		case "/":
			return func(i *item) {
				i.level = 1
			}
		default:
			log.Fatalf("invalid op [%v]", ss[0])
			return nil
		}
	}

	num, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Fatalf("not an int [%v]", ss[1])
	}
	switch ss[0] {
	case "+":
		return func(i *item) {
			i.level += num
		}
	case "-":
		return func(i *item) {
			i.level -= num
		}
	case "*":
		return func(i *item) {
			i.level *= num
		}
	case "/":
		return func(i *item) {
			i.level /= num
		}
	default:
		log.Fatalf("invalid op [%v]", ss[0])
		return nil
	}
}
