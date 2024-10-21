package main

import (
	"fmt"
	"strings"
	"time"
)

var zones = []string{"zone1", "zone2", "zone3", "zone4"}
var accessToken string

// Simulated token refresh logic
func refreshToken() string {
	return "newAccessToken"
}

func main() {
	// Infinite loop to refresh the access token every 30 minutes
	for {
		accessToken = refreshToken()

		// Collect results for the current interval
		var allZoneResults [][]Result

		// Process data for all zones
		for _, zone := range zones {
			entries := loadEntries() // Assuming entries are loaded from some source
			for i := 0; i < len(entries); i += 10 {
				end := i + 10
				if end > len(entries) {
					end = len(entries)
				}
				// Process each batch for the current zone
				results, err := ProcessZoneBatch(entries[i:end], zone, accessToken)
				if err != nil {
					fmt.Println("Error processing batch:", err)
					return
				}
				allZoneResults = append(allZoneResults, results)
				time.Sleep(30 * time.Second) // Wait for 30 seconds between batches
			}
		}

		// Write results to CSV
		outputFile := fmt.Sprintf("output_%s.csv", strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "_"))
		err := writeResultsToCSV(allZoneResults, outputFile)
		if err != nil {
			fmt.Println("Error writing to CSV:", err)
			return
		}

		fmt.Printf("Results written to %s\n", outputFile)

		// Sleep for 2 minutes before the next cycle
		fmt.Println("Waiting for 2 minutes before the next cycle...")
		time.Sleep(2 * time.Minute)
	}
}
