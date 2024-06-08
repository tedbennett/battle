package assert

import (
	"log"
	"testing"
)

func Assert(truthy bool, message string) {
	if !truthy {
		log.Fatalf(message)

	}
}

func TestAssert(t *testing.T, truthy bool, message string) {
	if !truthy {
		t.Fatalf(message)
	}
}
