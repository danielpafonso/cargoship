package configurations

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Service runs concatenation on files as process
type Service struct {
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
	Newline    bool   `yaml:"newline"`
}

// PachagerConfig structure describing configuration file
type PackagerConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"manifest"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Services []Service `yaml:"services"`
}

// ShipperReadConfig reads file configuration and return object/structure with shipper conconfiguration
func PackagerReadConfig(filepath string) (*PackagerConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config PackagerConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}

	// replace date format strings
	config.Log.Script = replaceDatePlaceholder(config.Log.Script)
	config.Log.Files = replaceDatePlaceholder(config.Log.Files)

	return &config, nil
}
