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

	file, err := os.Open("08/input.txt")
	//file, err := os.Open("08/test")
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
	log.Println("success: num visible trees", numVisibleTrees(parseGrid(bufio.NewScanner(file))))
}

func part2(file *os.File) {
	log.Println("success: highest score:", highestScenicScore(parseGrid(bufio.NewScanner(file))))
}

func parseGrid(scanner *bufio.Scanner) [][]int {
	result := make([][]int, 0)
	for scanner.Scan() {
		result = append(result, parseLine(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scan", err)
	}

	return result
}

func parseLine(str string) []int {
	result := make([]int, len(str))
	for i, c := range str {
		result[i] = intFrom(c)
	}
	return result
}

func intFrom(r rune) int {
	switch r {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	default:
		log.Fatalf("invalid rune [%v]", r)
		return 0
	}
}

func numVisibleTrees(grid [][]int) int {
	fromLeft := make([][]int, len(grid))
	for i, row := range grid {
		fromLeft[i] = make([]int, len(row))
		for j := 1; j < len(row); j++ {
			if row[j-1] > fromLeft[i][j-1] {
				fromLeft[i][j] = row[j-1]
				continue
			}
			fromLeft[i][j] = fromLeft[i][j-1]
		}
	}

	fromRight := make([][]int, len(grid))
	for i, row := range grid {
		fromRight[i] = make([]int, len(row))
		for j := len(row) - 2; j >= 0; j-- {
			if row[j+1] > fromRight[i][j+1] {
				fromRight[i][j] = row[j+1]
				continue
			}
			fromRight[i][j] = fromRight[i][j+1]
		}
	}

	fromUp := make([][]int, len(grid))
	fromUp[0] = make([]int, len(grid[0]))
	for i := 1; i < len(grid); i++ {
		fromUp[i] = make([]int, len(grid[i]))
		for j := range grid[i] {
			if fromUp[i-1][j] > grid[i-1][j] {
				fromUp[i][j] = fromUp[i-1][j]
				continue
			}
			fromUp[i][j] = grid[i-1][j]
		}
	}

	fromDown := make([][]int, len(grid))
	fromDown[len(fromDown)-1] = make([]int, len(grid[len(grid)-1]))
	for i := len(grid) - 2; i >= 0; i-- {
		fromDown[i] = make([]int, len(grid[i]))
		for j := range grid[i] {
			if fromDown[i+1][j] > grid[i+1][j] {
				fromDown[i][j] = fromDown[i+1][j]
				continue
			}
			fromDown[i][j] = grid[i+1][j]
		}
	}

	var result int
	for i, row := range grid {
		for j, tree := range row {
			if i == 0 || i == len(grid)-1 ||
				j == 0 || j == len(row)-1 ||
				tree > fromLeft[i][j] ||
				tree > fromRight[i][j] ||
				tree > fromUp[i][j] ||
				tree > fromDown[i][j] {
				result++
			}
		}
	}

	return result
}

func highestScenicScore(grid [][]int) int {
	type point struct {
		blockedAt *int // coordinate
		score     int
	}

	lefts := make([][]point, len(grid))
	for i, row := range grid {
		lefts[i] = make([]point, len(row))
		for j, tree := range row {
			if j == 0 {
				continue
			}
			score, curr := 1, ptrOf(j-1) // at least see one
			for curr != nil && grid[i][*curr] < tree {
				score += lefts[i][*curr].score
				curr = lefts[i][*curr].blockedAt
			}
			lefts[i][j] = point{
				blockedAt: curr,
				score:     score,
			}
		}
	}

	rights := make([][]point, len(grid))
	for i, row := range grid {
		rights[i] = make([]point, len(row))
		for j := len(row) - 1; j >= 0; j-- {
			if j == len(row)-1 {
				continue
			}
			score, curr := 1, ptrOf(j+1)
			for curr != nil && grid[i][*curr] < grid[i][j] {
				score += rights[i][*curr].score
				curr = rights[i][*curr].blockedAt
			}
			rights[i][j] = point{
				blockedAt: curr,
				score:     score,
			}
		}
	}

	ups := make([][]point, len(grid))
	for i, row := range grid {
		ups[i] = make([]point, len(row))
		if i == 0 {
			continue
		}
		for j, tree := range row {
			score, curr := 1, ptrOf(i-1)
			for curr != nil && grid[*curr][j] < tree {
				score += ups[*curr][j].score
				curr = ups[*curr][j].blockedAt
			}
			ups[i][j] = point{
				blockedAt: curr,
				score:     score,
			}
		}
	}

	downs := make([][]point, len(grid))
	for i := len(grid) - 1; i >= 0; i-- {
		downs[i] = make([]point, len(grid[i]))
		if i == len(grid)-1 {
			continue
		}
		for j, tree := range grid[i] {
			score, curr := 1, ptrOf(i+1)
			for curr != nil && grid[*curr][j] < tree {
				score += downs[*curr][j].score
				curr = downs[*curr][j].blockedAt
			}
			downs[i][j] = point{
				blockedAt: curr,
				score:     score,
			}
		}
	}

	var max int
	for i := 1; i < len(grid)-1; i++ {
		for j := 1; j < len(grid[i])-1; j++ {
			if curr := lefts[i][j].score * rights[i][j].score * ups[i][j].score * downs[i][j].score; curr > max {
				max = curr
			}
		}
	}
	return max
}

func ptrOf(i int) *int {
	return &i
}
