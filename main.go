package main

import (
	"fmt"
	"os"

	"gostarterapp/service"

	"github.com/spf13/cobra"
)

func main() {
	var env string
	var project string
	var bearerToken string

	var rootCmd = &cobra.Command{
		Use:   "gostarterapp",
		Short: "A CLI tool to fetch zones based on environment and project",
		Run: func(cmd *cobra.Command, args []string) {
			// Call GetZones function which uses validation to prompt for inputs and check validity
			env, project, bearerToken, zones, err := service.GetValidatedZones(env, project, bearerToken)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Print input details and zones
			fmt.Println("Environment:", env)
			fmt.Println("Project:", project)
			fmt.Println("Bearer Token:", bearerToken)
			fmt.Println("Zones:", zones)
		},
	}

	// Define flags
	rootCmd.Flags().StringVarP(&env, "env", "e", "", "Environment")
	rootCmd.Flags().StringVarP(&project, "project", "p", "", "Project")
	rootCmd.Flags().StringVarP(&bearerToken, "token", "t", "", "Bearer Token")

	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
