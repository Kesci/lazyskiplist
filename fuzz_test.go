package lazyskiplist

import (
	"testing"
)

func TestFuzz(t *testing.T) {
	// Fuzz([]byte{191, 0, 239, 23, 55, 55, 50, 48, 51, 57, 54, 53, 50, 127, 232, 161, 65, 184, 242})
	// Fuzz([]byte{239, 91, 37, 100, 47, 102, 105, 110, 100, 32, 108, 97, 121, 101, 114, 32, 37, 100, 191, 189, 23})
	Fuzz([]byte{239, 127, 0, 239, 127, 0, 0, 1, 249, 127, 40, 239, 127, 0, 0})
}