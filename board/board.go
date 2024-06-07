package board

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
	Team1Color = "afeeee"
	Team2Color = "eeafaf"
)

func teamToColor(team int) string {
	switch team {
	case Team1:
		return Team1Color
	case Team2:
		return Team2Color
	}
	return "000000"
}

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
