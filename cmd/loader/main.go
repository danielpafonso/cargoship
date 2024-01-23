package main

import (
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"
	"cargoship/internal/manifests"
)

var (
	scriptLogger logging.Logger
	filesLogger  logging.Logger
)

func compressFile(filename string, service configurations.LoaderServiceConfig) error {
	var dstPath string
	var archiver io.Writer

	if service.Archive == "gzip" {
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

func cleanFile(filename string, service configurations.LoaderServiceConfig) error {
	filepath := fmt.Sprintf("%s/%s", service.Src, filename)
	// clean files
	err := os.Remove(filepath)
	if err == nil {
		scriptLogger.LogInfo(fmt.Sprintf("Deleted file %s\n", filepath))
	}
	return err
}

func main() {
	// command line Flags
	var configFilepath string
	flag.StringVar(&configFilepath, "config", "loader_config.yaml", "Path to configuration yaml")
	flag.StringVar(&configFilepath, "c", "loader_config.yaml", "")
	flag.Usage = func() {
		fmt.Print(`Usage of shipper:
  -c, --config  path to configuation yaml 
  -h, --help    display this help message
`)
	}
	flag.Parse()

	// read script configuration file
	configs, err := configurations.LoaderReadConfig(configFilepath)
	if err != nil {
		log.Panic(err)
	}

	// start loggers
	scriptLogger.Init(configs.Log.Script, configs.Log2Console)
	filesLogger.Init(configs.Log.Files, configs.Log2Console)
	defer scriptLogger.Close()
	defer filesLogger.Close()

	// read load manifest
	manifest, err := manifests.LoaderReadTimes(configs.TimesPath)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}

	// defer update/write manifest file
	defer manifests.LoaderWriteTimes(&manifest, configs.TimesPath)

	// process services
	for _, service := range configs.Services {
		scriptLogger.LogInfo(fmt.Sprintf("Processing service: %s\n", service.Name))
		var processFile func(string, configurations.LoaderServiceConfig) error
		if service.Mode == "compress" {
			// set processFile function to compressFile
			processFile = compressFile
			files.CheckLocalFolder(service.Dst)
		} else if service.Mode == "cleaner" {
			// set processFile function to cleanFile
			processFile = cleanFile
		} else {
			scriptLogger.LogWarn(fmt.Sprintf("ERROR: Unkown mode, %s, on service %s.\n", service.Mode, service.Name))
			continue
		}

		// list local files
		files2Process, err := files.ListLocalDirectory(service.Src, service.Prefix, service.Extension)
		if err != nil {
			scriptLogger.LogError(err.Error())
			log.Panic(err)
		}
		// get last process file
		lastFile := manifests.LoaderGetTimes(&manifest, service.Name, service.Mode)
		files2Process = files.DateFilterLocalDirectory(files2Process, lastFile, service.MaxTime, service.Window)
		lastFileProcess := lastFile
		filesProcessed := 0
		for _, file := range files2Process {
			// if compress compress
			err = processFile(file.Name(), service)
			if err != nil {
				scriptLogger.LogError(err.Error())
			} else {
				filesProcessed += 1
			}
			//update
			lastFileProcess = file.ModTime()
		}

		scriptLogger.LogInfo(fmt.Sprintf("%d file(s) Processed", filesProcessed))
		// update last process file
		if lastFileProcess != lastFile {
			manifests.LoaderUpsertTimes(&manifest, service.Name, service.Mode, lastFileProcess)
		}
	}
}
