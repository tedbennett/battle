package board

import (
	"math/rand/v2"

	"github.com/tedbennett/battles/utils"
)

// Board represents which squares are taken by which side
type Board struct {
	Squares [][]int
	Teams   map[int]utils.Color
}

// These will be dynamic
const (
	Team1 int = iota
	Team2
)

const (
	Team1Color = "#afeeee"
	Team2Color = "#eeafaf"
)

func NewBoard(size int, element int, colors map[int]utils.Color) Board {
	squares := make([][]int, size)
	for i := range size {
		row := make([]int, size)
		for j := range size {
			row[j] = element
		}
		squares[i] = row
	}
	return Board{squares, colors}
}

// ========================================================
// Server side rendering of board
// ========================================================

func (b *Board) Colors() [][]string {
	colors := make([][]string, len(b.Squares))
	for rowIdx, row := range b.Squares {
		colorRow := make([]string, len(row))
		for colIdx, col := range row {
			colorRow[colIdx] = teamToColor(col)
		}
		colors[rowIdx] = colorRow
	}
	return colors
}

func teamToColor(team int) string {
	switch team {
	case Team1:
		return Team1Color
	case Team2:
		return Team2Color
	}
	return "000000"
}

// ========================================================
// Moving bar
// ========================================================

func NewBarBoard(size int, colors map[int]utils.Color) Board {
	b := NewBoard(size, Team1, colors)
	for i := range size {
		b.Squares[0][i] = Team2
	}
	return b
}

func (b *Board) TickBar() []Diff {
	diffs := make([]Diff, 0)
	for rowIdx, row := range b.Squares {
		if row[0] == Team2 {
			// Wipe prev
			for colIdx := range len(row) {
				b.Squares[rowIdx][colIdx] = Team1
				diffs = append(diffs, Diff{rowIdx, colIdx, Team1})
			}
			newRow := rowIdx + 1
			if rowIdx == len(row)-1 {
				newRow = 0
			}
			for colIdx := range len(row) {
				b.Squares[newRow][colIdx] = Team2
				diffs = append(diffs, Diff{newRow, colIdx, Team2})
			}
			break
		}
	}
	return diffs
}

// ========================================================
// Conways game of life
// ========================================================

func NewConwayBoard(size int, colors map[int]utils.Color) Board {
	b := NewBoard(size, Team1, colors)
	for range size {
		row, col := rand.IntN(size), rand.IntN(size)
		b.Squares[row][col] = Team2
	}
	return b
}

type Diff struct {
	Row, Col, Team int
}

func (b *Board) Tick() []Diff {
	diffs := make([]Diff, 0)
	for rowIdx, row := range b.Squares {
		for colIdx, col := range row {
			neighbours := b.getNeighbours(rowIdx, colIdx)
			if col == 1 {
				if neighbours > 3 || neighbours < 2 {
					b.Squares[rowIdx][colIdx] = Team1
					diffs = append(diffs, Diff{rowIdx, colIdx, Team1})
				}
			} else if neighbours == 3 {
				b.Squares[rowIdx][colIdx] = Team2
				diffs = append(diffs, Diff{rowIdx, colIdx, Team2})
			}
		}
	}
	return diffs
}

func (b *Board) getNeighbours(row int, col int) int {
	var count int = 0
	if row != 0 {
		count += b.Squares[row-1][col]
	}
	if row != len(b.Squares[0])-1 {
		count += b.Squares[row+1][col]
	}

	if col != 0 {
		count += b.Squares[row][col-1]
	}
	// Assuming square board
	if col != len(b.Squares[0])-1 {
		count += b.Squares[row][col+1]
	}
	return count
}
