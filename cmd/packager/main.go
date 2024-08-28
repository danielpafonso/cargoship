package main

import (
	"flag"
	"fmt"
	"log"

	"cargoship/internal/configurations"
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

	fmt.Printf("%+v\n", configs)
	for i := 0; i < len(configs.Services); i++ {
		configs.Services[i].Execute()
	}
}
