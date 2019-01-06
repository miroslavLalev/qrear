package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ConfigParser struct {
	configDir string
}

func NewConfigParser(configDir string) *ConfigParser {
	return &ConfigParser{configDir: configDir}
}

func (cp *ConfigParser) Parse(obj interface{}) error {
	files, err := ioutil.ReadDir(cp.configDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		f, err := os.Open(filepath.Join(cp.configDir, file.Name()))
		if err != nil {
			// TODO: log
			continue
		}
		err = json.NewDecoder(f).Decode(obj)
		if err != nil {
			// TODO: log
			continue
		}
	}
	return nil
}
