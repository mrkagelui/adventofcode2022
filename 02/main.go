package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("2/input.txt")
	if err != nil {
		log.Fatal("opening", err)
	}
	defer file.Close()

	var point int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) != 3 {
			log.Fatal("invalid line len:", text)
		}
		p, err := check(text[0], text[2])
		if err != nil {
			log.Fatal("invalid line:", text, " ", err)
		}
		point += p
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	log.Println("success", point)
}

func check(opponent, me uint8) (int, error) {
	switch opponent {
	case 'A':
		switch me {
		case 'X':
			return 3, nil
		case 'Y':
			return 4, nil
		case 'Z':
			return 8, nil
		default:
			return 0, fmt.Errorf("invalid me: %v", me)
		}
	case 'B':
		switch me {
		case 'X':
			return 1, nil
		case 'Y':
			return 5, nil
		case 'Z':
			return 9, nil
		default:
			return 0, fmt.Errorf("invalid me: %v", me)
		}
	case 'C':
		switch me {
		case 'X':
			return 2, nil
		case 'Y':
			return 6, nil
		case 'Z':
			return 7, nil
		default:
			return 0, fmt.Errorf("invalid me: %v", me)
		}
	default:
		return 0, fmt.Errorf("invalid opponent: %v", opponent)
	}
}
