package configurations

import (
	"os"

	"gopkg.in/yaml.v3"
)

// shipperFileConfig structure describing configuration file
type shipperFileConfig struct {
	Log2Console bool   `yaml:"log2console"`
	TimesPath   string `yaml:"manifest"`
	Log         struct {
		Script string `yaml:"script"`
		Files  string `yaml:"files"`
	} `yaml:"logging"`
	Ftps []struct {
		Name     string `yaml:"name"`
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
		User     string `yaml:"username"`
		Pass     string `yaml:"password"`
		Protocol string `yaml:"protocol"`
	} `yaml:"ftps"`
	Services []struct {
		Name      string   `yaml:"name"`
		Enable    bool     `yaml:"enable"`
		Ftp       []string `yaml:"ftpConfig"`
		Mode      string   `yaml:"mode"`
		Src       string   `yaml:"sourceFolder"`
		Dst       string   `yaml:"destinationFolder"`
		Prefix    string   `yaml:"filePrefix"`
		Extension string   `yaml:"fileExtension"`
		History   string   `yaml:"historyFolder"`
		MaxTime   int      `yaml:"maxTime"`
		Window    int      `yaml:"windowLimit"`
	} `yaml:"services"`
}

// ShipperService describes service configurations
type ShipperService struct {
	Name      string
	Enable    bool
	Mode      string
	Src       string
	Dst       string
	Prefix    string
	Extension string
	History   string
	MaxTime   int
	Window    int
}

// FtpConfig describes server connection configurations
type FtpConfig struct {
	Name     string
	Hostname string
	Port     int
	User     string
	Pass     string
	Protocol string
	Services []ShipperService
}

// ShipperConfig describes shipper run configuration
type ShipperConfig struct {
	Log2Console bool
	TimesPath   string
	Log         struct {
		Script string
		Files  string
	}
	Ftps []FtpConfig
}

// ShipperReadConfig reads file configuration and return object/structure with shipper conconfiguration
func ShipperReadConfig(filepath string) (*ShipperConfig, error) {
	// read file
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config shipperFileConfig
	// unmarshall it
	err = yaml.Unmarshal(fdata, &config)
	if err != nil {
		return nil, err
	}

	// process read config to match/duplicate service and server
	var newConfig = &ShipperConfig{
		Log2Console: config.Log2Console,
		TimesPath:   config.TimesPath,
		Log: struct {
			Script string
			Files  string
		}{
			Script: replaceDatePlaceholder(config.Log.Script),
			Files:  replaceDatePlaceholder(config.Log.Files),
		},
	}

	// create mapping service index
	ftpIndex := make(map[string]int, len(config.Ftps))
	for idx, ftp := range config.Ftps {
		ftpIndex[ftp.Name] = idx
		newConfig.Ftps = append(newConfig.Ftps, FtpConfig{
			Name:     ftp.Name,
			Hostname: ftp.Hostname,
			Port:     ftp.Port,
			User:     ftp.User,
			Pass:     ftp.Pass,
			Protocol: ftp.Protocol,
		})
	}

	// process serving matching to servers
	for _, service := range config.Services {
		match := ShipperService{
			Name:      service.Name,
			Enable:    service.Enable,
			Mode:      service.Mode,
			Src:       service.Src,
			Dst:       service.Dst,
			Prefix:    service.Prefix,
			Extension: service.Extension,
			History:   service.History,
			MaxTime:   service.MaxTime,
			Window:    service.Window,
		}
		for _, ftpName := range service.Ftp {
			if idx, ok := ftpIndex[ftpName]; ok {
				newConfig.Ftps[idx].Services = append(newConfig.Ftps[idx].Services, match)
			}
		}
	}
	return newConfig, nil
}
