package messages

import (
	"encoding/json"
	"fmt"

	"github.com/tedbennett/battles/assert"
)

const (
	INIT_MSG = iota
	BOARD_MSG
)

const VERSION int8 = 1

type Message struct {
	Type   int8
	Board  [][]int8
	Colors map[int]string
}

func NewInitMessage(colors map[int]string, board [][]int8) Message {
	return Message{Type: INIT_MSG, Colors: colors, Board: board}
}

func NewBoardMessage(board [][]int8) Message {
	return Message{Type: BOARD_MSG, Colors: nil, Board: board}
}

// Message format
// 1st byte: version
// 2nd byte: type
// 3+4th bytes: length
// 5-nth bytes: RLE encoded frame
func (m Message) Marshal() ([]byte, error) {
	buf := make([]byte, 0, len(m.Board))
	switch m.Type {
	case INIT_MSG:
		return json.Marshal(struct {
			Type  int8     `json:"type"`
			Board [][]int8 `json:"board"`
		}{
			Type:  m.Type,
			Board: m.Board,
		})

	case BOARD_MSG:
		buf = append(buf, byte(VERSION))
		buf = append(buf, byte(m.Type))
		buf = append(buf, []byte{0, 0}...)
		encoded := RleEncode(m.Board)
		Write16(buf, 2, len(encoded))
		buf = append(buf, encoded...)
		return buf, nil
	}
	return nil, fmt.Errorf("failed to marshall: invalid message type")
}
func Write16(buf []byte, offset, value int) {
	assert.Assert(len(buf) > offset+1, "you cannot write outside of the buffer")

	hi := (value & 0xFF00) >> 8
	lo := value & 0xFF
	buf[offset] = byte(hi)
	buf[offset+1] = byte(lo)
}
func RleEncode(board [][]int8) []byte {
	assert.Assert(len(board) > 0, "unable to encode an empty board")
	var current int8
	var count uint8
	buf := make([]byte, 0, len(board))
	for _, row := range board {
		for _, square := range row {
			if current == square && count < 255 {
				count += 1
			} else if current != square && count == 0 {
				// uninitialized
				current = square
				count += 1
			} else {
				buf = append(buf, byte(current))
				buf = append(buf, byte(count))
				count = 1
				current = square
			}
		}
	}

	buf = append(buf, byte(current))
	buf = append(buf, byte(count))
	return buf
}
