package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var ErrFileExists = os.ErrExist

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
		if err == ErrFileExists {
			LogInfo("Copier", fmt.Sprintf("%s already exists at %s", cp.src, cp.dst))
			return nil
		}

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
		err := os.MkdirAll(dstDir, os.ModeType)
		if err != nil {
			return err
		}

		return c.unsafeCopy(src, dst)
	}
}

func (c *Copier) unsafeCopy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// if the destination file exists, return an error
	//	if _, err := os.Stat(dst); err == nil {
	//		return ErrFileExists
	//	} else if errors.Is(err, os.ErrNotExist) {
	//		return ErrFileExists
	//	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
