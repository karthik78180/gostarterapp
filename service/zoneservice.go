package service

import (
	"fmt"

	"gostarterapp/models"
	"gostarterapp/validation"
)

func GetZones(env, project, bearerToken string) ([]string, error) {
	// Load config.json
	config, err := models.LoadConfig("configs/config.json")
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %v", err)
	}

	// Validate env and project
	if err := validation.ValidateEnv(env, config.Environments); err != nil {
		return nil, err
	}
	if err := validation.ValidateProject(env, project, config); err != nil {
		return nil, err
	}

	// Load zones.json
	zonesData, err := models.LoadZones("configs/zones.json")
	if err != nil {
		return nil, fmt.Errorf("Error loading zones: %v", err)
	}

	// Get the zones
	zones, err := zonesData.GetZones(env, project)
	if err != nil {
		return nil, err
	}

	return zones, nil
}
