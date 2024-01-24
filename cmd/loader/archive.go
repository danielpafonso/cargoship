package main

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"cargoship/internal/configurations"
)

func compressFile(filename string, service configurations.LoaderServiceConfig) error {
	var dstPath string
	var archiver io.Writer

	if service.Archive == "gz" {
		dstPath = fmt.Sprintf("%s/%s.gz", service.Dst, filename)
		// create archive file
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		archive := gzip.NewWriter(dstFile)
		defer archive.Close()
		archiver = archive
	} else if service.Archive == "zip" {
		dstPath = fmt.Sprintf("%s/%s.zip", service.Dst, filename)
		// create archive file
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		archive := zip.NewWriter(dstFile)
		defer archive.Close()
		archiver, err = archive.Create(filename)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown archive, %s", service.Archive)
	}

	// local file
	srcPath := fmt.Sprintf("%s/%s", service.Src, filename)
	localFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// archive file
	written, err := io.Copy(archiver, localFile)
	if err != nil {
		return err
	}
	// delete source file
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	filesLogger.LogInfo(fmt.Sprintf("Compress %s, written %d\n", dstPath, written))
	return nil
}

func uncompressFile(filename string, service configurations.LoaderServiceConfig) error {
	var dstPath string
	var archiveFile io.Reader

	srcPath := fmt.Sprintf("%s/%s", service.Src, filename)
	switch service.Archive {
	case "ungz":
		// read archive file
		fileReader, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer fileReader.Close()
		// get output path
		// EXT = .gz, len(EXT) = 3
		dstPath = fmt.Sprintf("%s/%s", service.Dst, filename[:len(filename)-3])
		// create gzip reader
		gzipReader, err := gzip.NewReader(fileReader)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		archiveFile = gzipReader

	case "unzip":
		// open zip file
		zipReader, err := zip.OpenReader(srcPath)
		if err != nil {
			return err
		}
		defer zipReader.Close()
		fileReader, err := zipReader.Open(filename[:len(filename)-4])
		if err != nil {
			return err
		}
		defer fileReader.Close()

		// get output path
		// EXT = .zip, len(EXT) = 4
		dstPath = fmt.Sprintf("%s/%s", service.Dst, filename[:len(filename)-4])
		archiveFile = fileReader

	default:
		return fmt.Errorf("unknown archive, %s", service.Archive)
	}

	// Create decompress file
	uncompress, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	// decompress file
	written, err := io.Copy(uncompress, archiveFile)
	if err != nil {
		return err
	}
	// delete source file
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	filesLogger.LogInfo(fmt.Sprintf("Decompress %s, written %d\n", dstPath, written))
	return nil
}
