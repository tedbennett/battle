package utils

import "github.com/tedbennett/battles/assert"

func Write16(buf []byte, offset, value int) {
	assert.Assert(len(buf) > offset+1, "buffer too small to write 2 bytes")

	hi := (value & 0xFF00) >> 8
	lo := value & 0xFF
	buf[offset] = byte(hi)
	buf[offset+1] = byte(lo)
}

func Read16(buf []byte, offset int) int {
	assert.Assert(len(buf) > offset+1, "buffer too small to read 2 bytes")

	hi := int(buf[offset])
	lo := int(buf[offset+1])
	return hi<<8 | lo
}
