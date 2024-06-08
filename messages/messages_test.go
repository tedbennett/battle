package messages_test

import (
	"reflect"
	"testing"

	"github.com/tedbennett/battles/assert"
	"github.com/tedbennett/battles/messages"
)

func TestRleEncode(t *testing.T) {
	board := [][]int8{{0, 0, 0, 0}, {1, 1, 1, 1}, {1, 0, 0, 0}}

	expected := []byte{byte(0), byte(4), byte(1), byte(5), byte(0), byte(3)}
	buf := messages.RleEncode(board)
	t.Log(expected)
	t.Log(buf)

	assert.TestAssert(t, reflect.DeepEqual(buf, expected), "failed to rle encode board")
}

func TestLongRleEncode(t *testing.T) {
	board := make([][]int8, 0)
	row := make([]int8, 16)
	for i := range 16 {
		row[i] = 1
	}
	for range 16 {
		board = append(board, row)
	}

	expected := []byte{byte(1), byte(255), byte(1), byte(1)}
	buf := messages.RleEncode(board)
	t.Log(expected)
	t.Log(buf)

	assert.TestAssert(t, reflect.DeepEqual(buf, expected), "failed to rle encode sequence of more than 255 repeated chars")
}
