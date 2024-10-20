package validation

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gostarterapp/models"
)

// Function to get environment, prompts if empty and shows available options
func GetEnv(env string, validEnvs []string) string {
	if env == "" {
		fmt.Printf("Enter environment from one of %v: ", validEnvs)
		reader := bufio.NewReader(os.Stdin)
		env, _ = reader.ReadString('\n')
		env = strings.TrimSpace(env)
	}
	return env
}

// Function to get project, prompts if empty and shows available options
func GetProject(project string, env string, config *models.Config) string {
	var availableProjects []string
	// Restrict projects to 'OD' and 'ODS' if the environment is 'dev'
	if env == "dev" {
		availableProjects = []string{"OD", "ODS"}
	} else {
		availableProjects = config.Projects
	}

	if project == "" {
		fmt.Printf("Enter project from one of %v: ", availableProjects)
		reader := bufio.NewReader(os.Stdin)
		project, _ = reader.ReadString('\n')
		project = strings.TrimSpace(project)
	}

	return project
}

// Function to get bearer token, prompts if empty
func GetBearerToken(bearerToken string) string {
	if bearerToken == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter bearer token: ")
		bearerToken, _ = reader.ReadString('\n')
		bearerToken = strings.TrimSpace(bearerToken)
	}
	return bearerToken
}

// ValidateEnv checks if the environment is valid
func ValidateEnv(env string, validEnvs []string) error {
	for _, e := range validEnvs {
		if env == e {
			return nil
		}
	}
	return fmt.Errorf("Invalid environment: '%s'", env)
}

// ValidateProject checks if the project is valid for the environment
func ValidateProject(env, project string, config *models.Config) error {
	// For 'dev', only 'OD' and 'ODS' are valid projects
	if env == "dev" {
		if project == "OD" || project == "ODS" {
			return nil
		}
		return fmt.Errorf("Project '%s' is not valid for environment 'dev'", project)
	}

	// For other environments, check if the project is valid
	for _, p := range config.Projects {
		if project == p {
			return nil
		}
	}
	return fmt.Errorf("Invalid project: '%s'", project)
}
