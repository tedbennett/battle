package board_test

import (
	"reflect"
	"testing"

	"github.com/tedbennett/battles/board"
)

func TestNewBoard(t *testing.T) {
	squares := [][]int{
		{board.Team1, board.Team1},
		{board.Team1, board.Team1},
	}
	expected := board.Board{Squares: squares}

	newBoard := board.NewBoard(2, board.Team1)

	Assert(t, reflect.DeepEqual(newBoard, expected), "failed to construct default board correctly")
}

func TestTeamToColor(t *testing.T) {
	squares := [][]int{
		{board.Team1, board.Team1},
		{board.Team1, board.Team2},
	}
	newBoard := &board.Board{Squares: squares}

	colors := newBoard.Colors()

	expected := [][]string{
		{board.Team1Color, board.Team1Color},
		{board.Team1Color, board.Team2Color},
	}

	Assert(t, reflect.DeepEqual(colors, expected), "failed to convert teams to colors")

}

func Assert(t *testing.T, truthy bool, message string) {
	if !truthy {
		t.Fatalf(message)
	}
}
