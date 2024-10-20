package models

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Environments []string `json:"environments"`
	Projects     []string `json:"projects"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
