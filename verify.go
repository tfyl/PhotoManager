package main

import (
	"os"
	"path/filepath"
	"strings"
)

func verifyFolder(src *string) error {
	// check if the input folder exists
	if _, err := os.Stat(*src); os.IsNotExist(err) {
		return err
	}

	// return full path
	abs, err := filepath.Abs(*src)
	if err != nil {
		return err
	}

	*src = abs
	return nil
}

func parseFileTypes(i string) ([]string, error) {
	iSplit := strings.Split(i, ",")
	// remove any empty strings/spaces
	for i, v := range iSplit {
		iSplit[i] = strings.TrimSpace(v)
	}

	return iSplit, nil
}
