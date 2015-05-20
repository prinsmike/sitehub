package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type JSONData map[string]interface{}

func validatePath(p string) bool {
	dir, err := os.Stat(p)
	if err != nil {
		return false
	}
	if !dir.IsDir() {
		return false
	}
	return true
}

func validateFile(fp string) bool {
	_, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return true
}

func host(rawHost string) string {
	return strings.TrimPrefix(strings.Split(rawHost, ":")[0], "www.")
}

func hostPath(rawHost string) (string, error) {
	hp := path.Join(config.workDir(), host(rawHost))
	if validatePath(hp) {
		return hp, nil
	}
	return "", errors.New(fmt.Sprintf("Could not find a valid path for host %s", hp))
}

func checkWorkDir() error {
	if !validatePath(config.workDir()) {
		err := os.MkdirAll(config.workDir(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadJSONFile(filePath string, data *map[string]interface{}) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, data)
	if err != nil {
		return err
	}
	return nil
}
