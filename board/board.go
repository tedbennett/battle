package board

import "math/rand/v2"

// Board represents which squares are taken by which side
type Board struct {
	Squares [][]int
}

// These will be dynamic
const (
	Team1 = iota
	Team2
)

const (
	Team1Color = "#afeeee"
	Team2Color = "#eeafaf"
)

func NewBoard(size int, element int) Board {
	squares := make([][]int, size)
	for i := range size {
		row := make([]int, size)
		for j := range size {
			row[j] = element
		}
		squares[i] = row
	}
	return Board{Squares: squares}
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
// Conways game of life
// ========================================================

func NewConwayBoard(size int) Board {
	b := NewBoard(size, Team1)
	for range size * 100 {
		row, col := rand.IntN(size), rand.IntN(size)
		b.Squares[row][col] = Team2
	}
	return b
}

func (b *Board) Tick() {
	for rowIdx, row := range b.Squares {
		for colIdx, col := range row {
			neighbours := b.getNeighbours(rowIdx, colIdx)
			if col == 1 {
				if neighbours > 3 || neighbours < 2 {
					b.Squares[rowIdx][colIdx] = 0
				}
			} else if neighbours == 3 {
				b.Squares[rowIdx][colIdx] = 1
			}
		}
	}
}

func (b *Board) getNeighbours(row int, col int) int {
	count := 0
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
