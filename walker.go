package main

import (
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

func walkerProgressBar(ch <-chan CounterType) {
	fileBar := progressbar.Default(-1)
	dirBar := progressbar.Default(-1)

	for cType := range ch {
		switch cType {
		case File:
			fileBar.Add(1)
		case Dir:
			dirBar.Add(1)
		}
	}

	dirBar.Finish()
	fileBar.Finish()
}

func Walker(dir string, fileTypesRegex *regexp.Regexp, recursive bool) ([]string, error) {
	ch := make(chan CounterType)
	defer close(ch)
	go walkerProgressBar(ch)

	files, dirs, err := walker(dir, fileTypesRegex, recursive, ch)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		newFiles, newDirs, err := walker(dir, fileTypesRegex, recursive, ch)
		if err != nil {
			return nil, err
		}
		files = append(files, newFiles...)
		dirs = append(dirs, newDirs...)
	}

	return files, nil
}

func walker(dir string, fileTypesRegex *regexp.Regexp, recursive bool, counterCh chan CounterType) ([]string, []string, error) {
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
				counterCh <- Dir
			}
		} else {
			if fileTypesRegex.MatchString(file.Name()) {
				files = append(files, filepath.Join(dir, file.Name()))
				counterCh <- File
			}
		}
	}

	return files, dirs, nil
}
