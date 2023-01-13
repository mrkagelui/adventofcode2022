package main

import "fmt"

type bottomedGrid [][]spot

func withRocksWithBottom(lines []line) bottomedGrid {
	const sourceCol = 500

	var maxRow int
	for _, l := range lines {
		if l.start.row > maxRow {
			maxRow = l.start.row
		}
		if l.end.row > maxRow {
			maxRow = l.end.row
		}
	}

	height := maxRow + 3 // + 1 to make room for the lowest rock; + 2 for bottom
	omitted := sourceCol - (height - 12)
	width := height*2 - 1 // because the sand can't span past this before source is blocked

	g := make(bottomedGrid, height)
	for i := range g {
		g[i] = make([]spot, width)
	}

	for _, l := range lines {
		if l.start.row == l.end.row { // horizontal line
			var s, e int
			if l.start.col <= l.end.col {
				s, e = l.start.col-omitted, l.end.col-omitted
			} else {
				s, e = l.end.col-omitted, l.start.col-omitted
			}
			for i := s; i <= e; i++ {
				if i < 0 || i > width-1 {
					continue
				}
				g[l.start.row][i].isRock = true
			}
			continue
		}
		if l.start.col == l.end.col { // vertical line
			if l.start.col-omitted < 0 || l.start.col-omitted > width-1 {
				continue
			}

			var s, e int
			if l.start.row <= l.end.row {
				s, e = l.start.row, l.end.row
			} else {
				s, e = l.end.row, l.start.row
			}
			for i := s; i <= e; i++ {
				g[i][l.start.col-omitted].isRock = true
			}
		}
	}

	// set bottom (actually not necessary, purely aesthetic)
	for i := 0; i < width; i++ {
		g[height-1][i].isRock = true
	}

	return g
}

func (g bottomedGrid) show() {
	for i := 0; i < len(g)-1; i++ {
		fmt.Print(" ")
	}
	fmt.Println("|")
	for _, spots := range g {
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

func (g bottomedGrid) maxSand() int {
	fmt.Println("before")
	g.show()
	sourceCol := len(g) - 1

	var i int
	for !g[0][sourceCol].isSand {
		g.dropOneSand()
		i++
	}

	fmt.Println("after")
	g.show()
	return i
}

func (g bottomedGrid) dropOneSand() {
	sourceCol := len(g) - 1

	sand := point{row: 0, col: sourceCol}
	for sand.row < len(g)-2 { // don't need to touch the second last row
		// no bound check needed, the width and height is calculated

		// try below
		if g[sand.row+1][sand.col].canGo() {
			sand.row++
			continue
		}

		// try bottom left
		if g[sand.row+1][sand.col-1].canGo() {
			sand.row++
			sand.col--
			continue
		}

		// try bottom right
		if g[sand.row+1][sand.col+1].canGo() {
			sand.row++
			sand.col++
			continue
		}

		g[sand.row][sand.col].isSand = true
		return
	}
	g[sand.row][sand.col].isSand = true
}
