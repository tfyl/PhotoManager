package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

func Walker(dir string, fileTypesRegex *regexp.Regexp, recursive bool) ([]string, error) {
	files, dirs, err := walker(dir, fileTypesRegex, recursive)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		newFiles, newDirs, err := walker(dir, fileTypesRegex, recursive)
		if err != nil {
			return nil, err
		}
		files = append(files, newFiles...)
		dirs = append(dirs, newDirs...)
	}

	return files, nil
}

func walker(dir string, fileTypesRegex *regexp.Regexp, recursive bool) ([]string, []string, error) {
	var files []string
	var dirs []string

	allFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range allFiles {
		if file.IsDir() {
			if recursive {
				dirs = append(dirs, filepath.Join(dir, file.Name()))
			}
		} else {
			if fileTypesRegex.MatchString(file.Name()) {
				files = append(files, filepath.Join(dir, file.Name()))
			}
		}
	}

	return files, dirs, nil
}
