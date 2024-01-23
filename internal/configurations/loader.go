package configurations

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoaderServiceConfig stores service configuration
type LoaderServiceConfig struct {
	Name      string `yaml:"name"`
	Mode      string `yaml:"mode"`
	Src       string `yaml:"sourceFolder"`
	Dst       string `yaml:"destinationFolder"`
	Prefix    string `yaml:"filePrefix"`
	Extension string `yaml:"fileExtension"`
	Archive   string `yaml:"archive"`
	MaxTime   int    `yaml:"maxTime"`
	Window    int    `yaml:"windowLimit"`
}

// LoaderConfig structure describing loader run configuration, obtain directly form file
type LoaderConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"manifest"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Services []LoaderServiceConfig `yaml:"services"`
}

// LoaderReadConfig reads file configuration and return object/structure with shipper conconfiguration
func LoaderReadConfig(filepath string) (*LoaderConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config LoaderConfig
	// unmarshal data
	err = yaml.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}

	// replace date format strings
	config.Log.Script = replaceDatePlaceholder(config.Log.Script)
	config.Log.Files = replaceDatePlaceholder(config.Log.Files)

	return &config, nil
}
