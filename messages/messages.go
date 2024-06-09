package messages

import (
	"bytes"
	"fmt"

	"github.com/tedbennett/battles/assert"
	"github.com/tedbennett/battles/board"
	"github.com/tedbennett/battles/utils"
)

const (
	INIT_MSG = iota
	BOARD_MSG
	PARTIAL_MSG
)

var MESSAGE_MAP = map[int]string{
	INIT_MSG:    "Init",
	BOARD_MSG:   "Board",
	PARTIAL_MSG: "Partial",
}

const VERSION = 1

type Payload interface {
	Type() byte
	Write(buf []byte) (int, error)
	Debug(buf []byte) string
}

type Message struct {
	buf     []byte
	Payload Payload
}

func NewMessage(size int, p Payload) *Message {
	m := &Message{
		buf:     make([]byte, size, size),
		Payload: p,
	}
	m.Write()
	return m
}

/*
Header:
1	 2	  3	   4
+--------+--------+--------+--------+
| vers.  |  type  |  payload len    |
+--------+--------+--------+--------+
*/
func (m *Message) PackHeader(out []byte, offset int) {
	assert.Assert(len(out) >= 4, "header buffer too short")
	out[offset] = byte(VERSION)
	out[offset+1] = byte(m.Payload.Type())
}

func (m *Message) Write() error {
	m.PackHeader(m.buf, 0)
	m.Payload.Write(m.buf[2:])
	return nil
}

func (m *Message) Bytes() []byte {
	return m.buf
}

func (m *Message) Debug() {
	fmt.Printf("Version: %b, Type: %s\n", m.buf[0], MESSAGE_MAP[int(m.buf[1])])
	fmt.Printf("Payload: %s", m.Payload.Debug(m.buf[2:]))
}

type InitMessage struct {
	colors map[int]utils.Color
	b      *board.Board
}

func NewInitMessage(c map[int]utils.Color, b *board.Board) *InitMessage {
	return &InitMessage{c, b}
}

func (i *InitMessage) Type() byte {
	return byte(INIT_MSG)
}

const BYTES_PER_COLOR = 4

func (i *InitMessage) Write(buf []byte) (int, error) {
	// 4 bytes per team:color pair
	n := len(i.colors) * BYTES_PER_COLOR
	encoded := RleEncode(i.b.Squares)
	assert.Assert(len(buf) >= n+len(encoded)+4, "buffer too small to write InitMessage")
	utils.Write16(buf, 0, n)
	offset := 2

	for team, color := range i.colors {
		buf[offset] = byte(team)
		color.Write(buf[offset+1:])
		offset += BYTES_PER_COLOR
	}

	utils.Write16(buf, n+2, len(encoded))
	copy(buf[n+4:], encoded)

	return n + len(encoded) + 4, nil
}

func (i *InitMessage) Debug(buf []byte) string {
	var str bytes.Buffer
	length := int(utils.Read16(buf, 0))
	str.WriteString(fmt.Sprintf("Len: %d\n", length))
	for i := 0; i < length; i += BYTES_PER_COLOR {
		offset := i + 2
		str.WriteString(fmt.Sprintf("Team %d: %s\n", int(buf[offset]), utils.DebugColor(buf[offset+1:offset+4])))
	}

	bLength := int(utils.Read16(buf, length+2))
	str.WriteString(fmt.Sprintf("Len: %d\n", bLength))
	// for i := 0; i < bLength; i += 2 {
	// 	offset := length + i + 4
	// 	str.WriteString(fmt.Sprintf("Count %d: Char %d\n", int(buf[offset]), int(buf[offset+1])))
	// }
	return str.String()
}

type BoardMessage struct {
	board *board.Board
}

func NewBoardMessage(b *board.Board) *BoardMessage {
	return &BoardMessage{b}
}

func (i *BoardMessage) Type() byte {
	return byte(BOARD_MSG)
}

func (i *BoardMessage) Write(buf []byte) (int, error) {
	encoded := RleEncode(i.board.Squares)
	assert.Assert(len(buf) >= len(encoded)+2, "buffer too small to write BoardMessage")
	utils.Write16(buf, 0, len(encoded))
	copy(buf[2:], encoded)
	return len(encoded) + 2, nil
}

func (i *BoardMessage) Debug(buf []byte) string {
	var str bytes.Buffer
	length := int(utils.Read16(buf, 0))
	str.WriteString(fmt.Sprintf("Len: %d\n", length))
	for i := 0; i < length; i += 2 {
		offset := i + 2
		str.WriteString(fmt.Sprintf("Count %d: Char %d\n", int(buf[offset]), int(buf[offset+1])))
	}
	return str.String()
}

type PartialMessage struct {
	diffs []board.Diff
}

func NewPartialMessage(diffs []board.Diff) *PartialMessage {
	return &PartialMessage{diffs}
}

func (i *PartialMessage) Type() byte {
	return byte(PARTIAL_MSG)
}

const BYTES_PER_DIFF = 3

func (i *PartialMessage) Write(buf []byte) (int, error) {
	// Format of Row, Col, Team
	n := len(i.diffs) * BYTES_PER_DIFF
	assert.Assert(len(buf) >= n+2, "buffer too small to write PartialMessage")
	utils.Write16(buf, 0, n)
	for i, diff := range i.diffs {
		offset := (i * BYTES_PER_DIFF) + 2
		buf[offset] = byte(diff.Row)
		buf[offset+1] = byte(diff.Col)
		buf[offset+2] = byte(diff.Team)
	}
	return n + 2, nil
}

func (i *PartialMessage) Debug(buf []byte) string {
	var str bytes.Buffer
	length := int(utils.Read16(buf, 0))
	str.WriteString(fmt.Sprintf("Len: %d\n", length))
	for i := 0; i < length-2; i += BYTES_PER_DIFF {
		offset := i + 2
		str.WriteString(fmt.Sprintf("Row %d, Col %d, Team %d\n", int(buf[offset]), int(buf[offset+1]), int(buf[offset+2])))
	}
	return str.String()
}

func RleEncode(board [][]int) []byte {
	assert.Assert(len(board) > 0, "unable to encode an empty board")
	var current int
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
				buf = append(buf, byte(count))
				buf = append(buf, byte(current))
				count = 1
				current = square
			}
		}
	}

	buf = append(buf, byte(count))
	buf = append(buf, byte(current))
	return buf
}
