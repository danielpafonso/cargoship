package transport

import (
	"cargoship/internal/logging"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/pkg/sftp"
)

// checkSshFolder check remote folder existence, creating folder if doesn't exist
func checkSshFolder(client *sftp.Client, folder string, logger logging.Logger) {
	_, err := client.Stat(folder)

	if err != nil {
		// folder doesn't exists, create
		logger.LogInfo(fmt.Sprintf("Create remote folder %s\n", folder))
		client.MkdirAll(folder)
	}
}

// listSshDirectory list remote files present on a folder, filtering files by filename prefix and extension
func listSshDirectory(client *sftp.Client, folder, prefix, extension string) ([]fs.FileInfo, error) {
	var outputList []fs.FileInfo

	files, err := client.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), extension) {
			outputList = append(outputList, file)
		}
	}
	return outputList, nil
}

// dateFilterSshDirectory filter remote file list, returns files that are after lastTime and before limit
func dateFilterSshDirectory(entries []fs.FileInfo, lastTime time.Time, maxTime int, limit int) []fs.FileInfo {
	var outputList []fs.FileInfo

	filesLimit := time.Now().UTC().Add(time.Minute * time.Duration(limit*-1))

	// sort files by modification time
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ModTime().Before(entries[j].ModTime())
	})

	for _, entry := range entries {
		if entry.ModTime().After(lastTime) && entry.ModTime().Before(filesLimit) {
			if len(outputList) == 0 {
				// update file limit with max time
				maxLimit := entry.ModTime().Add(time.Minute * time.Duration(maxTime))
				if maxLimit.Before(filesLimit) {
					filesLimit = maxLimit
				}
			}
			outputList = append(outputList, entry)
		}
		// cut for loop if files (entry) are after the file limit time
		if entry.ModTime().After(filesLimit) {
			break
		}
	}
	return outputList
}
