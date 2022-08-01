package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gammazero/workerpool"
	"log"
	"os"
	"time"
)

type Logging int64

const (
	Normal = iota
	Debug
)

var DebugLevel = Normal

func verifyFlags(inputFolder *string, outputFolder *string, fileType *string, recursive *bool, threads *int) ([]string, error) {
	LogInfo("PhotoManager", "Initialising...")
	flag.Parse()

	err := verifyFolder(inputFolder)
	if err != nil {
		return nil, FormatError("InputFolder", err)
	}

	err = verifyFolder(outputFolder)
	if err != nil {
		return nil, FormatError("OutputFolder", err)
	}

	fileTypeArr, err := parseFileTypes(*fileType)
	if err != nil {
		return nil, FormatError("FileType", err)
	}

	LogInfo("Settings", "Input folder: "+*inputFolder)
	LogInfo("Settings", fmt.Sprintf("File type: %v", fileTypeArr))
	LogInfo("Settings", "Output folder: "+*outputFolder)
	LogInfo("Settings", "Recursive: "+fmt.Sprintf("%t", *recursive))
	return fileTypeArr, nil
}

func main() {
	inputFolder := flag.String("i", "./input", "Input folder")
	fileType := flag.String("t", "jpeg,jpg,raf", "Type of file to process (enter extensions separated by commas)")
	outputFolder := flag.String("o", "./output", "Output folder (Directory the folders will be created in)")
	recursive := flag.Bool("r", true, "recursive: Includes any files found in the input subfolders")
	threads := flag.Int("threads", 5, "Number of threads to use")

	fileTypeArr, err := verifyFlags(inputFolder, outputFolder, fileType, recursive, threads)
	if err != nil {
		LogError("PhotoManager-Flags", err.Error())
		os.Exit(1)
	}

	fileTypeRe := FileTypeToRe(fileTypeArr)

	files, err := Walker(*inputFolder, fileTypeRe, *recursive)
	if err != nil {
		LogError("PhotoManager-Walker", err.Error())
		os.Exit(1)
	}

	LogInfo("PhotoManager", fmt.Sprintf("Found %d files", len(files)))

	p := &Parser{}

	copies, err := p.Parse(files, *outputFolder)
	if err != nil {
		return
	}

	progressCh := make(chan CounterType)
	defer close(progressCh)
	go copierProgressBar(progressCh, len(copies))

	c := &Copier{}
	wp := workerpool.New(*threads)

	for _, cp := range copies {
		cp := cp
		wp.Submit(func() {
			if DebugLevel == Debug {
				log.Printf("Copying %s to %s", cp.src, cp.dst)
			}

			err := c.Copy(cp)
			var fileErr ErrFileExists
			if err != nil && errors.As(err, &fileErr) {
				if DebugLevel == Debug {
					LogError("Copier", fmt.Sprintf("%s already exists at %s", fileErr.src, fileErr.dst))
				}
			} else if err != nil {
				LogError("Copier", err.Error())
			} else if DebugLevel == Debug {
				LogInfo("Copier", fmt.Sprintf("%s to %s", cp.src, cp.dst))
			}
			progressCh <- File
		})
	}

	wp.StopWait()
	LogSuccess("PhotoManager", "Finished")
	time.Sleep(time.Second * 10)
}
