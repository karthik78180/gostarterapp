package service

import (
	"fmt"

	"gostarterapp/models"
	"gostarterapp/validation"
)

func GetValidatedZones(env, project, bearerToken string) (string, string, string, []string, error) {
	// Load config.json
	config, err := models.LoadConfig("configs/config.json")
	if err != nil {
		return "", "", "", nil, fmt.Errorf("Error loading config: %v", err)
	}

	// Load valid environments
	validEnvs := config.Environments

	// Validate env and project using validation package
	env = validation.GetEnv(env, validEnvs) // Ensure env is obtained or prompted
	if err := validation.ValidateEnv(env, validEnvs); err != nil {
		return "", "", "", nil, err
	}

	project = validation.GetProject(project, env, config) // Ensure project is obtained or prompted
	if err := validation.ValidateProject(env, project, config); err != nil {
		return "", "", "", nil, err
	}

	// Prompt for bearer token if missing
	bearerToken = validation.GetBearerToken(bearerToken)

	// Load zones.json
	zonesData, err := models.LoadZones("configs/zones.json")
	if err != nil {
		return "", "", "", nil, fmt.Errorf("Error loading zones: %v", err)
	}

	// Get the zones
	zones, err := zonesData.GetZones(env, project)
	if err != nil {
		return "", "", "", nil, err
	}

	return env, project, bearerToken, zones, nil
}
