package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"time"
)

func parserProgressBar(ch <-chan CounterType, len int) {
	fileBar := progressbar.Default(int64(len))

	for cType := range ch {
		switch cType {
		case File:
			fileBar.Add(1)
		}
	}

	fileBar.Finish()
}

type Parser struct{}

func (p *Parser) Parse(srcs []string, dstFolder string) ([]*Copy, error) {
	ch := make(chan CounterType)
	defer close(ch)
	go parserProgressBar(ch, len(srcs))

	var copies []*Copy

	for _, src := range srcs {
		dst, err := p.parseFile(src)
		if err != nil {
			return nil, err
		}
		dst = filepath.Join(dstFolder, dst)

		copies = append(copies, NewCopy(src, dst))
		ch <- File
	}

	return copies, nil
}

func (p *Parser) parseFile(src string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return "", err
	}

	bTime, err := x.DateTime()
	if err != nil {
		return "", err
	}

	toDir := p.timeToDirName(bTime)

	return filepath.Join(toDir, filepath.Base(src)), nil
}

func (p *Parser) timeToDirName(t time.Time) string {
	return t.Format("2006-01-02")
}
