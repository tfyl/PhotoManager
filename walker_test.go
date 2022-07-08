package main

import (
	"testing"
)

func TestWalker(t *testing.T) {
	re := FileTypeToRe([]string{"raf"})

	files, err := Walker(`path/to/folder`, re, true)
	if err != nil {
		t.Error(err)
	}

	t.Logf("files: %v", files)
}
