package utils_test

import (
	"testing"

	"github.com/tedbennett/battles/assert"
	"github.com/tedbennett/battles/utils"
)

func TestColorFromString(t *testing.T) {
	_, err := utils.SafeColorFromString("#ff00ff")
	assert.TestAssert(t, err == nil, "failed to parse #ff00ff")

	_, err = utils.SafeColorFromString("ff00ff")
	assert.TestAssert(t, err != nil, "failed to raise error for missing #")

	_, err = utils.SafeColorFromString("fg00ff")
	assert.TestAssert(t, err != nil, "failed to raise error for invalid hex color string")
}

func TestColorWrite(t *testing.T) {
	color := utils.ColorFromString("#ff00ff")
	buf := make([]byte, 3, 3)

	n, err := color.Write(buf)
	assert.TestAssert(t, err == nil, "failed to write color to buffer")
	assert.TestAssert(t, n == 3, "did not write expected number of bytes for color")
	assert.TestAssert(t, int(buf[0]) == 255, "invalid byte written for color")
	assert.TestAssert(t, int(buf[1]) == 0, "invalid byte written for color")
	assert.TestAssert(t, int(buf[2]) == 255, "invalid byte written for color")
}
