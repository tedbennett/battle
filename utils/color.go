package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/tedbennett/battles/assert"
)

type Color struct {
	r uint8
	g uint8
	b uint8
}

func ColorFromString(str string) Color {
	// Assume all ASCII
	color, err := SafeColorFromString(str)
	assert.Assert(err == nil, "unable to parse hex string")
	return color
}

func SafeColorFromString(str string) (Color, error) {
	// Assume all ASCII
	colors, err := hex.DecodeString(str[1:])
	if err != nil {
		return Color{}, err
	}
	if len(colors) != 3 {
		return Color{}, fmt.Errorf("incorrect number of hex components")
	}
	assert.Assert(len(colors) == 3, "unable to parse hex string")
	return Color{
		r: colors[0],
		g: colors[1],
		b: colors[2],
	}, nil
}

func (c *Color) Write(buf []byte) (int, error) {
	buf[0] = byte(c.r)
	buf[1] = byte(c.g)
	buf[2] = byte(c.b)
	return 3, nil
}

func DebugColor(buf []byte) string {
	return fmt.Sprintf("r: %d, g: %d, b: %d", int(buf[0]), int(buf[1]), int(buf[2]))
}
