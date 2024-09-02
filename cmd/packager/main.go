package main

import (
	"flag"
	"fmt"
	"log"

	"cargoship/internal/configurations"
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
	fmt.Printf("%+v\n", configs)
	for _, service := range configs.Services {
		service.Execute()
	}
}
