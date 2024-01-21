package transport

import (
	"fmt"
	"io"
	"log"
	"os"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"
	"cargoship/internal/manifests"

	"github.com/pkg/sftp"
)

// sftpDownload a single file from remote folder to a local folder
func sftpDownload(client *sftp.Client, source, destination string, logger logging.Logger) error {

	srcFile, err := client.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create local writer
	fmt.Println(destination)
	localFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer localFile.Close()

	sizeWritten, err := io.Copy(localFile, srcFile)
	if err != nil {
		return err
	}
	//logger.LogInfo(fmt.Sprintf("Donwloaded file %s (size %d), written %d\n", entry.Name, entry.Size, sizeWritten))
	logger.LogInfo(fmt.Sprintf("Donwloaded file %s, written %d\n", source, sizeWritten))

	return nil
}

// SftpDownloadFiles downloads files described in configurations from remote folder to local folder
func SftpDownloadFiles(
	serverName string,
	client *sftp.Client,
	service configurations.ShipperService,
	times *[]manifests.ShipperManifest,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {
	// check local folder
	files.CheckLocalFolder(service.Dst)
	// check remote history folder
	if service.History != "" {
		checkSshFolder(client, service.History, scriptLogger)
	}
	// get last file time
	fileTime := manifests.ShipperGetTimes(*times, serverName, service.Mode, service.Name)
	// list files in directory
	entries, err := listSshDirectory(client, service.Src, service.Prefix, service.Extension)
	if err != nil {
		log.Panic(err)
	}

	entries = dateFilterSshDirectory(entries, fileTime, service.MaxTime, service.Window)
	// check if there are any files to download
	if len(entries) == 0 {
		scriptLogger.LogInfo("No files to download")
		return
	}
	// donwload files
	lastFileTime := fileTime

	for _, entry := range entries {
		src := fmt.Sprintf("%s/%s", service.Src, entry.Name())
		dst := fmt.Sprintf("%s/%s", service.Dst, entry.Name())
		err := sftpDownload(client, src, dst, filesLogger)
		if err != nil {
			scriptLogger.LogWarn(err.Error())
			break
		}
		if service.History != "" {
			err := client.Rename(src, fmt.Sprintf("%s/%s", service.History, entry.Name()))
			if err != nil {
				scriptLogger.LogWarn(err.Error())
			}
			scriptLogger.LogInfo(fmt.Sprintf("Moved file %s to history folder %s\n", entry.Name(), service.History))
		}
		// update
		lastFileTime = entry.ModTime()
		fmt.Println(lastFileTime)
	}

	// update last downloaded time

	if lastFileTime != fileTime {
		manifests.ShipperUpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}

// sftpUpload a single file from local folder to a remote folder
func sftpUpload(client *sftp.Client, source, destination string, logger logging.Logger) error {
	// local reader
	localReader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer localReader.Close()

	// upload
	remoteFile, err := client.Create(destination)
	if err != nil {
		return err
	}
	remoteSize, err := io.Copy(remoteFile, localReader)
	if err != nil {
		return err
	}
	//logger.LogInfo(fmt.Sprintf("Uploaded file %s (size %d), written %d\n", entry.Name(), entry.Size(), remoteSize))
	logger.LogInfo(fmt.Sprintf("Uploaded file %s, written %d\n", source, remoteSize))

	return nil
}

// SftpUploadFiles uploads files described in configurations from local folder to remote folder
func SftpUploadFiles(
	serverName string,
	client *sftp.Client,
	service configurations.ShipperService,
	times *[]manifests.ShipperManifest,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {

	scriptLogger.LogInfo(fmt.Sprintf("Processing %s: %s\n", service.Mode, service.Name))
	// check folders
	checkSshFolder(client, service.Dst, scriptLogger)
	if service.History != "" {
		files.CheckLocalFolder(service.History)
	}
	// get last file time
	fileTime := manifests.ShipperGetTimes(*times, serverName, service.Mode, service.Name)

	// list files in directory
	entries, err := files.ListLocalDirectory(service.Src, service.Prefix, service.Extension)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}

	entries = files.DateFilterLocalDirectory(entries, fileTime, service.MaxTime, service.Window)
	// check if there are any files to upload
	if len(entries) == 0 {
		scriptLogger.LogInfo("No files to upload")
		return
	}
	// upload files
	lastFileTime := fileTime
	for _, entry := range entries {
		src := fmt.Sprintf("%s/%s", service.Src, entry.Name())
		dst := fmt.Sprintf("%s/%s", service.Dst, entry.Name())

		err := sftpUpload(client, src, dst, filesLogger)
		if err != nil {
			scriptLogger.LogWarn(err.Error())
			break
		}
		if service.History != "" {
			os.Rename(
				src,
				fmt.Sprintf("%s/%s", service.History, entry.Name()),
			)
			scriptLogger.LogInfo(fmt.Sprintf("Moved file %s to history folder %s\n", entry.Name(), service.History))
		}
		// updatuploade
		lastFileTime = entry.ModTime()
	}
	if lastFileTime != fileTime {
		manifests.ShipperUpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}
