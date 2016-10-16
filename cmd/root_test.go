package cmd

import (
	"os"
	"testing"
)

func TestFileStitch(t *testing.T) {
	chunkSize = 4
	chunkFile("test.txt")
	os.Remove("test.txt")
	stitchOriginal("test.txt")
}
