package cmd

import (
	"testing"
)

func TestFileStitch(t *testing.T) {
	chunkSize = 4
	chunkFile("test.txt")
}
