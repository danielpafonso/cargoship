package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"
)

func ConcatFiles(filesProcess []os.FileInfo, lastProcessedTime time.Time, serviceConfig configurations.Service, logger logging.Logger) (time.Time, int, error) {
	filesProcessed := 0

	// prep work
	files.CheckLocalFolder(serviceConfig.Dst)
	if serviceConfig.History != "" {
		files.CheckLocalFolder(serviceConfig.History)
	}
	// create output file in tmp folder
	tmpOutput, err := os.CreateTemp("", "packagerprocesstmp")
	if err != nil {
		return lastProcessedTime, 0, err
	}

	// process files
	for _, file := range filesProcess {
		scriptLogger.LogDebug(file.Name())
		// create reader
		reader, err := os.Open(path.Join(serviceConfig.Src, file.Name()))
		if err != nil {
			logger.LogError(err.Error())
			continue
		}
		defer reader.Close()
		// copy file contents
		_, err = io.Copy(tmpOutput, reader)
		if err != nil {
			logger.LogError(err.Error())
			os.Remove(tmpOutput.Name())
		}
		// Add newline if defined
		if serviceConfig.Newline {
			tmpOutput.WriteString("\n")
		}
		// update processed stat
		filesProcessed += 1
		lastProcessedTime = file.ModTime()

		// mode source to history folder
		if serviceConfig.History != "" {
			err = files.MoveFile(path.Join(serviceConfig.Src, file.Name()), path.Join(serviceConfig.History, file.Name()))
			if err != nil {
				return lastProcessedTime, 0, err
			}
		}
	}
	// generate output filename
	outFileName := strings.Clone(serviceConfig.Output)
	outFileName = strings.Replace(outFileName, "{date}", time.Now().UTC().Format(serviceConfig.DateFormat), 1)
	outFileName = strings.Replace(outFileName, "{files}", fmt.Sprint(filesProcessed), 1)

	// move output file to destination
	err = files.MoveFile(tmpOutput.Name(), path.Join(serviceConfig.Dst, outFileName))
	if err != nil {
		return lastProcessedTime, 0, err
	}
	return lastProcessedTime, filesProcessed, nil
}
