package manifests

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// packagerFileManifest structure ofr packager file manifest
type packagerFileManifest struct {
	Service string `yaml:"service"`
	Mode    string `yaml:"mode"`
	Time    string `yaml:"time"`
}

// PackagerManifest store files times (manifest) for packager services
type PackagerManifest struct {
	Service string
	Mode    string
	Time    time.Time
}

// PackagerReadTimes read manifest files
func PackagerReadTimes(filepath string) ([]PackagerManifest, error) {
	// create empty readConfig
	var readConfig []packagerFileManifest

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
	var config []PackagerManifest
	for _, serviceService := range readConfig {
		timestamp, _ := time.Parse(time.RFC3339, serviceService.Time)
		config = append(config, PackagerManifest{
			Service: serviceService.Service,
			Mode:    serviceService.Mode,
			Time:    timestamp,
		})
	}

	return config, nil
}

// PackagerGetTimes search manifest and returns date from last process file, or "empty" date if service isn't stored
func PackagerGetTimes(manifests *[]PackagerManifest, service, mode string) time.Time {
	for _, elm := range *manifests {
		if elm.Service == service && elm.Mode == mode {
			return elm.Time
		}
	}
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

// PackagerUpsertTimes insert/update serice manifest
func PackagerUpsertTimes(manifests *[]PackagerManifest, service, mode string, times time.Time) {
	for idx, manifest := range *manifests {
		if manifest.Service == service && manifest.Mode == mode {
			(*manifests)[idx].Time = times
			return
		}
	}
	// new ftp times
	*manifests = append(*manifests, PackagerManifest{Service: service, Mode: mode, Time: times})
}

// PackagerWriteTimes write manifest to file
func PackagerWriteTimes(ftpTimes *[]PackagerManifest, filepath string) error {
	fileData := make([]packagerFileManifest, len(*ftpTimes))
	for idx, stamp := range *ftpTimes {
		fileData[idx] = packagerFileManifest{
			Service: stamp.Service,
			Mode:    stamp.Mode,
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
