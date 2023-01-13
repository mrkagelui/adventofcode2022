package main

import "fmt"

type spot struct {
	isRock, isSand bool
}

func (s spot) canGo() bool {
	return !s.isSand && !s.isRock
}

// I know it's more "range-tolerant" to use some discrete data representation, but this is more visualized
type grid struct {
	horizontalRetraction int
	m                    [][]spot
	highestAirSpotRow    map[int]int // key: col, value: highest row of air spot before hitting a rock
}

func withRocks(lines []line) grid {
	minCol := lines[0].start.col
	var maxCol, maxRow int
	for _, l := range lines {
		if l.start.col < minCol {
			minCol = l.start.col
		}
		if l.start.col > maxCol {
			maxCol = l.start.col
		}
		if l.start.row > maxRow {
			maxRow = l.start.row
		}

		if l.end.col < minCol {
			minCol = l.end.col
		}
		if l.end.row > maxRow {
			maxRow = l.end.row
		}
		if l.end.col > maxCol {
			maxCol = l.end.col
		}
	}

	width, height := maxCol-minCol+1, maxRow+1
	matrix := make([][]spot, height)
	for i := range matrix {
		matrix[i] = make([]spot, width)
	}
	highestAirSpots := make(map[int]int, width)

	for _, l := range lines {
		if l.start.row == l.end.row { // horizontal line
			var s, e int
			if l.start.col <= l.end.col {
				s, e = l.start.col-minCol, l.end.col-minCol
			} else {
				s, e = l.end.col-minCol, l.start.col-minCol
			}
			for i := s; i <= e; i++ {
				matrix[l.start.row][i] = spot{isRock: true}
				highestSoFar, ok := highestAirSpots[i]
				if !ok || l.start.row+1 > highestSoFar {
					highestAirSpots[i] = l.start.row + 1 // this can go below the whole matrix
				}
			}
			continue
		}
		if l.start.col == l.end.col { // vertical line
			var s, e int
			if l.start.row <= l.end.row {
				s, e = l.start.row, l.end.row
			} else {
				s, e = l.end.row, l.start.row
			}
			for i := s; i <= e; i++ {
				matrix[i][l.start.col-minCol] = spot{isRock: true}
			}
			highestSoFar, ok := highestAirSpots[l.start.col]
			if !ok || e+1 > highestSoFar {
				highestAirSpots[l.start.col] = e + 1 // this can go below the whole matrix
			}
		}
	}

	return grid{
		horizontalRetraction: minCol,
		m:                    matrix,
		highestAirSpotRow:    highestAirSpots,
	}
}

func (g grid) show() {
	for i := 0; i < 500-g.horizontalRetraction; i++ {
		fmt.Print(" ")
	}
	fmt.Println("|")
	for _, spots := range g.m {
		for _, s := range spots {
			var c byte
			switch {
			case s.isSand:
				c = 'o'
			case s.isRock:
				c = '#'
			default:
				c = '.'
			}
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
}

func (g grid) maxSand() int {
	fmt.Println("before")
	g.show()
	var i int
	for g.dropOneSand() {
		i++
	}
	fmt.Println("after")
	g.show()
	return i
}

func (g grid) dropOneSand() (stopped bool) {
	const initCol = 500
	width, height := len(g.m[0]), len(g.m)

	sandPos := point{row: 0, col: initCol - g.horizontalRetraction}
	for sandPos.row < height-1 {
		if highestAir := g.highestAirSpotRow[sandPos.col]; sandPos.row >= highestAir {
			// falls through if no rock below
			return false
		}

		// try below (no need to check bound because the outer for loop guarantees that)
		if g.m[sandPos.row+1][sandPos.col].canGo() {
			sandPos.row++
			continue
		}

		// try bottom left
		if sandPos.col-1 < 0 { // fall out
			return false
		}
		if g.m[sandPos.row+1][sandPos.col-1].canGo() {
			sandPos.row++
			sandPos.col--
			continue
		}

		// try bottom right
		if sandPos.col+1 > width-1 { // fall out
			return false
		}
		if g.m[sandPos.row+1][sandPos.col+1].canGo() {
			sandPos.row++
			sandPos.col++
			continue
		}

		g.m[sandPos.row][sandPos.col].isSand = true
		return true
	}
	return false
}
