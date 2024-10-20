package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ZonesData map[string]map[string]struct {
	Zones []string `json:"zones"`
}

func LoadZones(filename string) (ZonesData, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var zonesData ZonesData
	err = json.Unmarshal(data, &zonesData)
	if err != nil {
		return nil, err
	}
	return zonesData, nil
}

func (zd ZonesData) GetZones(env, project string) ([]string, error) {
	if projects, ok := zd[env]; ok {
		if projectData, ok := projects[project]; ok {
			return projectData.Zones, nil
		} else {
			return nil, fmt.Errorf("Project '%s' not found in environment '%s'", project, env)
		}
	} else {
		return nil, fmt.Errorf("Environment '%s' not found", env)
	}
}
