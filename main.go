package main

import (
	"flag"
	"fmt"
	"github.com/gammazero/workerpool"
	"os"
	"time"
)

func verifyFlags(inputFolder *string, outputFolder *string, fileType *string, recursive *bool) ([]string, error) {
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
	recursive := flag.Bool("r", false, "recursive: Includes any files found in the input subfolders")

	fileTypeArr, err := verifyFlags(inputFolder, outputFolder, fileType, recursive)
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

	c := &Copier{}
	wp := workerpool.New(10)

	for _, cp := range copies {
		cp := cp
		wp.Submit(func() {
			err := c.Copy(cp)
			if err != nil {
				LogError("PhotoManager-Copier", err.Error())
			}
			LogInfo("Copier", fmt.Sprintf("%s to %s", cp.src, cp.dst))
		})
	}

	wp.StopWait()
	LogSuccess("PhotoManager", "Finished")
	time.Sleep(time.Second * 10)
}
