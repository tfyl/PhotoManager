package main

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type Parser struct{}

func (p *Parser) Parse(srcs []string, dstFolder string) ([]*Copy, error) {
	var copies []*Copy

	for _, src := range srcs {
		dst, err := p.parseFile(src)
		if err != nil {
			return nil, err
		}
		dst = filepath.Join(dstFolder, dst)

		copies = append(copies, NewCopy(src, dst))
	}

	return copies, nil
}

func (p *Parser) parseFile(src string) (string, error) {
	// gets file metadata
	fInfo, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	fData := fInfo.Sys().(*syscall.Win32FileAttributeData)

	eTime := earliestTime(fData.CreationTime.Nanoseconds(), fData.LastWriteTime.Nanoseconds())
	pTime := time.Unix(0, eTime)

	toDir := p.timeToDirName(pTime)

	return filepath.Join(toDir, fInfo.Name()), nil
}

func earliestTime(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (p *Parser) timeToDirName(t time.Time) string {
	return t.Format("2006-01-02")
}
