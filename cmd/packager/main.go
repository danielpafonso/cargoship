package main

import (
	"flag"
	"fmt"
	"log"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"
	"cargoship/internal/manifests"
)

var (
	scriptLogger logging.Logger
	filesLogger  logging.Logger
)

func main() {
	// command line Flags
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "packager_config.yaml", "Path to configuration yaml")
	flag.StringVar(&configFilePath, "c", "packager_config.yaml", "")
	flag.Usage = func() {
		fmt.Print(`Usage of packager:
  -c, --config  path to configuation yaml 
  -h, --help    display this help message
`)
	}
	flag.Parse()

	// read scripy configuration file
	configs, err := configurations.PackagerReadConfig(configFilePath)
	if err != nil {
		log.Panic(err)
	}

	// start loggers
	scriptLogger.Init(configs.Log.Script, configs.Log2Console)
	filesLogger.Init(configs.Log.Files, configs.Log2Console)
	defer scriptLogger.Close()
	defer filesLogger.Close()

	// read time state
	times, err := manifests.PackagerReadTimes(configs.TimesPath)
	if err != nil {
		scriptLogger.LogError(err.Error())
		panic(err)
	}
	// defer update/write time state file
	defer manifests.PackagerWriteTimes(&times, configs.TimesPath)

	// Process files"
	for _, service := range configs.Services {
		scriptLogger.LogInfo(fmt.Sprintf("Processing serice %s\n", service.Name))
		// set execute function

		// list local files
		files2Process, err := files.ListLocalDirectory(service.Name, service.Prefix, service.Mode)
		if err != nil {
			scriptLogger.LogError(err.Error())
			continue
		}
		// get last process file
		lastFileTime := manifests.PackagerGetTimes(&times, service.Name, service.Mode)
		// update list of local files
		files2Process = files.DateFilterLocalDirectory(files2Process, lastFileTime, service.MaxTime, service.Window)

		// process files
		lastFileProcess := lastFileTime
		filesProcessed := 0

	}
}
