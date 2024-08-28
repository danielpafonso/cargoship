package configurations

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// pachagerFileConfig structure describing configuration file
type packagerFileConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"manifest"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Services []struct {
		Name       string `yaml:"name"`
		Mode       string `yaml:"mode"`
		Src        string `yaml:"sourceFolder"`
		Prefix     string `yaml:"filePrefix"`
		Extension  string `yaml:"fileExtension"`
		Dst        string `yaml:"destinationFolder"`
		Output     string `yaml:"destinationFile"`
		DateFormat string `yaml:"destinationDateFormat"`
		History    string `yaml:"historyFolder"`
		MaxTime    int    `yaml:"maxTime"`
		Window     int    `yaml:"windowLimit"`
	} `yaml:"services"`
}

// PackagerService enables the execution of configured packager processing
type PackagerService interface {
	Execute() error
}

// ConcatService runs concatenation on files as process
type ConcatService struct {
	Name       string
	Mode       string
	Src        string
	Prefix     string
	Extension  string
	Dst        string
	Output     string
	DateFormat string
	History    string
	MaxTime    int
	Window     int
}

func (srv *ConcatService) Execute() error {
	fmt.Printf("execute %v\n", srv.Mode)
	return nil
}

// CommandService runs the configured command as process
type CommandService struct {
}

// ShipperConfig describes shipper run configuration
type PackagerConfig struct {
	Log2Console bool
	TimesPath   string
	Log         struct {
		Script string
		Files  string
	}
	Services []PackagerService
}

// ShipperReadConfig reads file configuration and return object/structure with shipper conconfiguration
func PackagerReadConfig(filepath string) (*PackagerConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config packagerFileConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}

	// process read config to match/duplicate service and server
	var newConfig = &PackagerConfig{
		Log2Console: config.Log2Console,
		TimesPath:   config.TimesPath,
		Log: struct {
			Script string
			Files  string
		}{
			Script: replaceDatePlaceholder(config.Log.Script),
			Files:  replaceDatePlaceholder(config.Log.Files),
		},
		Services: make([]PackagerService, 0),
	}

	// process serving matching to servers
	for _, service := range config.Services {
		newConfig.Services = append(newConfig.Services, &ConcatService{
			Name:       service.Name,
			Mode:       service.Mode,
			Src:        service.Src,
			Prefix:     service.Prefix,
			Extension:  service.Extension,
			Dst:        service.Dst,
			Output:     service.Output,
			DateFormat: service.DateFormat,
			History:    service.History,
			MaxTime:    service.MaxTime,
			Window:     service.Window,
		})
	}
	return newConfig, nil
}
