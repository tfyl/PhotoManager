package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func copierProgressBar(ch <-chan CounterType, len int) {
	fileBar := progressbar.Default(int64(len))

	for cType := range ch {
		switch cType {
		case File:
			fileBar.Add(1)
		}
	}

	fileBar.Finish()
}

type ErrFileExists struct {
	src string
	dst string
}

func NewErrFileExists(src, dst string) ErrFileExists {
	return ErrFileExists{src, dst}
}

func (e ErrFileExists) Error() string {
	return fmt.Sprintf("%s already exists at %s", e.src, e.dst)
}

type Copy struct {
	src string
	dst string
}

func NewCopy(src, dst string) *Copy {
	return &Copy{src, dst}
}

type Copier struct {
	// directories that are known to exist
	dirs sync.Map
}

func (c *Copier) Copy(cp *Copy) error {
	err := c.copy(cp.src, cp.dst)
	if err != nil {
		return err
	}

	return nil
}

func (c *Copier) copy(src, dst string) error {
	dstDir := filepath.Dir(dst)
	if _, ok := c.dirs.Load(dstDir); ok {
		return c.unsafeCopy(src, dst)
	} else {
		c.dirs.Store(dstDir, true)
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return FormatError("Copier-MkdirAll", err)
		}

		return c.unsafeCopy(src, dst)
	}
}

func (c *Copier) unsafeCopy(src, dst string) error {
	fDst, destErr := os.Stat(dst)
	if destErr == nil {
		fSrc, err := os.Stat(src)
		if err != nil {
			return FormatError("Copier-Stat-Nil", err)
		}

		if fSrc.Size() == fDst.Size() {
			return NewErrFileExists(src, dst)
		}

		if DebugLevel == Debug {
			LogInfo("Copier-Overwrite", fmt.Sprintf("%s already exists at %s, overwriting", src, dst))
		}
	}

	source, err := os.Open(src)
	if err != nil {
		return FormatError("Copier-Open", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return FormatError("Copier-Create", err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return FormatError("Copier-Copy", err)
	}

	return nil
}
