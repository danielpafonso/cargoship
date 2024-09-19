package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"
)

func CommandFiles(filesProcess []os.FileInfo, lastProcessedTime time.Time, serviceConfig configurations.Service, logger logging.Logger) (time.Time, int, error) {
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
	// defer os.Remove(tmpOutput.Name())

	// process files
	for _, file := range filesProcess {
		scriptLogger.LogDebug(file.Name())

		// prep command
		var cmdLine string
		if strings.Contains(serviceConfig.Command, "{file}") {
			cmdLine = strings.ReplaceAll(serviceConfig.Command, "{file}", path.Join(serviceConfig.Src, file.Name()))
		} else {
			cmdLine = fmt.Sprintf("%s %s", serviceConfig.Command, path.Join(serviceConfig.Src, file.Name()))
		}

		// execute Command
		cmdArray := strings.SplitN(cmdLine, " ", 2)
		cmd := exec.Command(cmdArray[0], cmdArray[1])
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			logger.LogError(err.Error())
			os.Remove(tmpOutput.Name())
		}
		// Add output to file
		tmpOutput.Write(cmdOutput)
		// Add newline if defined
		if serviceConfig.Newline {
			tmpOutput.WriteString("\n")
		}

		// update processed stat
		filesProcessed += 1
		lastProcessedTime = file.ModTime()

		// mode source to history folder
		if serviceConfig.History != "" {
			os.Rename(path.Join(serviceConfig.Src, file.Name()), path.Join(serviceConfig.History, file.Name()))
		}
	}
	// generate output filename
	outFileName := strings.Clone(serviceConfig.Output)
	if strings.Contains(outFileName, "{date}") {
		outFileName = strings.Replace(outFileName, "{date}", time.Now().UTC().Format(serviceConfig.DateFormat), 1)
	}
	if strings.Contains(outFileName, "{files}") {
		outFileName = strings.Replace(outFileName, "{files}", fmt.Sprint(filesProcessed), 1)
	}

	// move output file to destination
	os.Rename(tmpOutput.Name(), path.Join(serviceConfig.Dst, outFileName))

	return lastProcessedTime, filesProcessed, nil
}
