package validation

import (
	"fmt"
	"gostarterapp/models"
)

func ValidateEnv(env string, validEnvs []string) error {
	for _, e := range validEnvs {
		if env == e {
			return nil
		}
	}
	return fmt.Errorf("Invalid environment: '%s'", env)
}

func ValidateProject(env, project string, config *models.Config) error {
	// For 'dev', only 'OD' and 'ODS' are valid projects
	if env == "dev" {
		if project == "OD" || project == "ODS" {
			return nil
		}
		return fmt.Errorf("Project '%s' is not valid for environment 'dev'", project)
	}

	// For other environments, check if project is valid
	for _, p := range config.Projects {
		if project == p {
			return nil
		}
	}
	return fmt.Errorf("Invalid project: '%s'", project)
}
