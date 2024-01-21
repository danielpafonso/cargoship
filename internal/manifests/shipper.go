package manifests

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Shipper structures
type shipperFileManifest struct {
	Ftp     string `yaml:"ftp"`
	Mode    string `yaml:"mode"`
	Service string `yaml:"service"`
	Time    string `yaml:"time"`
}

// ShipperManifest store files times (manifest) for export and import services
type ShipperManifest struct {
	Ftp     string
	Mode    string
	Service string
	Time    time.Time
}

// ShipperReadTimes read manifest files
func ShipperReadTimes(filepath string) ([]ShipperManifest, error) {
	// create empty readConfig
	var readConfig []shipperFileManifest

	// check if file exits
	if _, err := os.Stat(filepath); err == nil {
		//read file
		fdata, err := os.ReadFile((filepath))
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(fdata, &readConfig)
		if err != nil {
			return nil, err
		}
	}

	// convert string into time object
	var config []ShipperManifest
	for _, serviceService := range readConfig {
		timestamp, _ := time.Parse(time.RFC3339, serviceService.Time)
		config = append(config, ShipperManifest{
			Ftp:     serviceService.Ftp,
			Mode:    serviceService.Mode,
			Service: serviceService.Service,
			Time:    timestamp,
		})
	}

	return config, nil
}

// ShipperGetTimes search manifest and returns date from last process file, or "empty" date if service isn't stored
func ShipperGetTimes(ftpTimes []ShipperManifest, server string, mode string, service string) time.Time {
	for _, elm := range ftpTimes {
		if elm.Ftp == server && elm.Service == service && elm.Mode == mode {
			return elm.Time
		}
	}
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

// ShipperUpsertTimes insert/update serice manifest
func ShipperUpsertTimes(config *[]ShipperManifest, ftp string, mode string, service string, times time.Time) {
	for idx, filetimes := range *config {
		if filetimes.Ftp == ftp && filetimes.Mode == mode && filetimes.Service == service {
			(*config)[idx].Time = times
			return
		}
	}
	// new ftp times
	*config = append(*config, ShipperManifest{Ftp: ftp, Mode: mode, Service: service, Time: times})
}

// ShipperWriteTimes writes files manivest to file
func ShipperWriteTimes(ftpTimes *[]ShipperManifest, filepath string) error {
	fileData := make([]shipperFileManifest, len(*ftpTimes))
	for idx, stamp := range *ftpTimes {
		fileData[idx] = shipperFileManifest{
			Ftp:     stamp.Ftp,
			Mode:    stamp.Mode,
			Service: stamp.Service,
			Time:    stamp.Time.Format(time.RFC3339),
		}
	}

	// marshall times -> struct to string
	data, err := yaml.Marshal(fileData)
	if err != nil {
		return err
	}
	// write to file
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
