package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type configuration struct {
	Sinks []sink `json:"sinks"`
}

type sink struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// loadConfig parses and validates the configuration file
func loadConfig(fpath string) (configuration, error) {
	buf, err := os.ReadFile(fpath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return configuration{},
				fmt.Errorf(
					"config file %q: %s. Run 'lanterna init' to create a configuration file",
					fpath, err)
		}
		return configuration{}, fmt.Errorf("config file %q: %s", fpath, err)
	}

	var config configuration
	if err := json.Unmarshal(buf, &config); err != nil {
		return configuration{}, err
	}

	if err := validateConfig(config); err != nil {
		return configuration{}, fmt.Errorf("config file %q: %s", fpath, err)
	}

	return config, nil
}

// TODO actually validate something!
func validateConfig(config configuration) error {
	if len(config.Sinks) == 0 {
		return fmt.Errorf("  empty configuration")
	}

	// FIXME validate should NOT know type names. Should invoke a registration function or similar.
	// FIXME should collect all errors instead of stopping at the first one

	for _, sink := range config.Sinks {
		if sink.Type != "gchat" {
			return fmt.Errorf("sink %q has unsupported type %q", sink.Name, sink.Type)
		}
	}

	return nil
}
