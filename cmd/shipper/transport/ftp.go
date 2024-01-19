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

	"github.com/jlaffaye/ftp"
)

// download a single file from remote folder to a local folder
func download(conn *ftp.ServerConn, destination string, entry *ftp.Entry, logger logging.Logger) error {
	remoteReader, err := conn.Retr(entry.Name)
	if err != nil {
		return err
	}
	defer remoteReader.Close()

	// create local writer
	localWriter, err := os.OpenFile(
		destination,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer localWriter.Close()

	sizeWritten, err := io.Copy(localWriter, remoteReader)
	if err != nil {
		return err
	}
	logger.LogInfo(fmt.Sprintf("Donwloaded file %s (size %d), written %d\n", entry.Name, entry.Size, sizeWritten))

	return nil
}

// DownloadFiles downloads files described in configurations from remote folder to local folder
func DownloadFiles(
	serverName string,
	ftpConn *ftp.ServerConn,
	service configurations.ShipperService,
	times *[]manifests.ShipperManifest,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {

	// check folder
	files.CheckLocalFolder(service.Dst)
	// check remote hostory folder
	if service.History != "" {
		checkRemoteFolder(ftpConn, service.History, scriptLogger)
	}

	// get last file time
	fileTime := manifests.ShipperGetTimes(*times, serverName, service.Mode, service.Name)
	// list files in directory
	entries, err := listRemoteDirectory(ftpConn, service.Src, service.Prefix, service.Extension)
	if err != nil {
		log.Panic(err)
	}

	entries = dateFilterRemoteDirectory(entries, fileTime, service.MaxTime, service.Window)
	// check if there are any files to download
	if len(entries) == 0 {
		scriptLogger.LogInfo("No files to download")
		return
	}
	// donwload files
	lastFileTime := fileTime
	err = ftpConn.ChangeDir(service.Src)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}
	for _, entry := range entries {
		err := download(ftpConn, fmt.Sprintf("%s/%s", service.Dst, entry.Name), entry, filesLogger)
		if err != nil {
			scriptLogger.LogWarn(err.Error())
			break
		}
		if service.History != "" {
			err := ftpConn.Rename(entry.Name, fmt.Sprintf("%s/%s", service.History, entry.Name))
			if err != nil {
				scriptLogger.LogWarn(err.Error())
			}
			scriptLogger.LogInfo(fmt.Sprintf("Moved file %s to history folder %s\n", entry.Name, service.History))
		}
		// update
		lastFileTime = entry.Time
	}
	// update last downloaded time
	if lastFileTime != fileTime {
		manifests.ShipperUpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}

// upload a single file from local folder to a remote folder
func upload(conn *ftp.ServerConn, source string, entry os.FileInfo, logger logging.Logger) error {
	// local reader
	localReader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer localReader.Close()

	// upload
	err = conn.Stor(entry.Name(), localReader)
	if err != nil {
		return err
	}
	remoteSize, _ := conn.FileSize(entry.Name())
	logger.LogInfo(fmt.Sprintf("Uploaded file %s (size %d), written %d\n", entry.Name(), entry.Size(), remoteSize))

	return nil
}

// UploadFiles uploads files described in configurations from local folder to remote folder
func UploadFiles(
	serverName string,
	ftpConn *ftp.ServerConn,
	service configurations.ShipperService,
	times *[]manifests.ShipperManifest,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {

	scriptLogger.LogInfo(fmt.Sprintf("Processing %s: %s\n", service.Mode, service.Name))
	// check folders
	checkRemoteFolder(ftpConn, service.Dst, scriptLogger)
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
	err = ftpConn.ChangeDir(service.Dst)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}
	for _, entry := range entries {
		err := upload(ftpConn, fmt.Sprintf("%s/%s", service.Src, entry.Name()), entry, filesLogger)
		if err != nil {
			scriptLogger.LogWarn(err.Error())
			break
		}
		if service.History != "" {
			os.Rename(
				fmt.Sprintf("%s/%s", service.Src, entry.Name()),
				fmt.Sprintf("%s/%s", service.History, entry.Name()),
			)
			scriptLogger.LogInfo(fmt.Sprintf("Moved file %s to history folder %s\n", entry.Name(), service.History))
		}
		// update
		lastFileTime = entry.ModTime()
	}
	if lastFileTime != fileTime {
		manifests.ShipperUpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}
