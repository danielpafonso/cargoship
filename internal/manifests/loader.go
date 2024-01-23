package manifests

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// loaderFileManifest structure ofr loader file manifest
type loaderFileManifest struct {
	Service string `yaml:"service"`
	Mode    string `yaml:"mode"`
	Time    string `yaml:"time"`
}

// LoaderManifest store files times (manifest) for loader services
type LoaderManifest struct {
	Service string
	Mode    string
	Time    time.Time
}

// LoaderReadTimes read manifest files
func LoaderReadTimes(filepath string) ([]LoaderManifest, error) {
	// create empty readConfig
	var readConfig []loaderFileManifest

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
	var config []LoaderManifest
	for _, serviceService := range readConfig {
		timestamp, _ := time.Parse(time.RFC3339, serviceService.Time)
		config = append(config, LoaderManifest{
			Service: serviceService.Service,
			Mode:    serviceService.Mode,
			Time:    timestamp,
		})
	}

	return config, nil
}

// LoaderGetTimes search manifest and returns date from last process file, or "empty" date if service isn't stored
func LoaderGetTimes(manifests *[]LoaderManifest, service, mode string) time.Time {
	for _, elm := range *manifests {
		if elm.Service == service && elm.Mode == mode {
			return elm.Time
		}
	}
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
}

// LoaderUpsertTimes insert/update serice manifest
func LoaderUpsertTimes(manifests *[]LoaderManifest, service, mode string, times time.Time) {
	for idx, manifest := range *manifests {
		if manifest.Service == service && manifest.Mode == mode {
			(*manifests)[idx].Time = times
			return
		}
	}
	// new ftp times
	*manifests = append(*manifests, LoaderManifest{Service: service, Mode: mode, Time: times})
}

// LoaderWriteTimes write manifest to file
func LoaderWriteTimes(ftpTimes *[]LoaderManifest, filepath string) error {
	fileData := make([]loaderFileManifest, len(*ftpTimes))
	for idx, stamp := range *ftpTimes {
		fileData[idx] = loaderFileManifest{
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
