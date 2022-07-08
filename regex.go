package main

import "regexp"

func FileTypeToRe(fileTypes []string) *regexp.Regexp {
	reStr := ".*((?i)"
	for _, fileType := range fileTypes {
		reStr += fileType + "|"
	}
	reStr = reStr[:len(reStr)-1]
	reStr += ")$"
	return regexp.MustCompile(reStr)
}
