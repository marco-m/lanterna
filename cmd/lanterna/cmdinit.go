package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

func cmdInit(cfg Args) error {
	if err := doInit(cfg.ConfigPath); err != nil {
		return fmt.Errorf("init: %s", err)
	}
	fmt.Printf("init: created %s\n", cfg.ConfigPath)
	return nil
}

func doInit(fpath string) error {
	if err := os.MkdirAll(path.Dir(fpath), 0700); err != nil {
		return err
	}

	_, err := os.Stat(fpath)
	if err == nil {
		return fmt.Errorf("file already exists: %s", fpath)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	config := configuration{
		Sinks: []sink{
			{
				Name: "sink name, suggested: room name",
				Type: "sink type",
				URL:  "sink webhook",
			},
		},
	}

	buf, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if _, err := file.Write(buf); err != nil {
		return err
	}

	return nil
}
