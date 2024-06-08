package messages

import (
	"encoding/json"
	"fmt"
)

const (
	INIT_MSG = iota
	BOARD_MSG
)

type Message struct {
	Type   int
	Board  [][]int
	Colors map[int]string
}

func NewInitMessage(colors map[int]string, board [][]int) Message {
	return Message{Type: INIT_MSG, Colors: colors, Board: board}
}

func NewBoardMessage(board [][]int) Message {
	return Message{Type: BOARD_MSG, Colors: nil, Board: board}
}

func (m Message) MarshalJSON() ([]byte, error) {
	switch m.Type {
	case INIT_MSG:
		return json.Marshal(struct {
			Type   int            `json:"type"`
			Board  [][]int        `json:"board"`
			Colors map[int]string `json:"colors"`
		}{
			Type:   m.Type,
			Board:  m.Board,
			Colors: m.Colors,
		})
	case BOARD_MSG:
		return json.Marshal(struct {
			Type  int     `json:"type"`
			Board [][]int `json:"board"`
		}{
			Type:  m.Type,
			Board: m.Board,
		})
	}
	return nil, fmt.Errorf("failed to marshall: invalid message type")
}
