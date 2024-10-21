package main

import (
	"fmt"
	"time"
)

var zones = []string{"zone1", "zone2", "zone3", "zone4"}
var accessToken string

// Simulated token refresh logic
func refreshToken() string {
	return "newAccessToken"
}

func main() {
	outputFile := "output_results.csv" // Single output file for all cycles

	// Infinite loop to refresh the access token every 30 minutes
	for {
		accessToken = refreshToken()

		// Step 1: Load all entries
		entries := loadEntries()

		// Step 2: Create batches of 10 entries
		var allZoneResults [][]Result
		for i := 0; i < len(entries); i += 10 {
			end := i + 10
			if end > len(entries) {
				end = len(entries)
			}
			batch := entries[i:end]

			// Step 3: Loop over all zones for the current batch of 10 entries
			for _, zone := range zones {
				results, err := ProcessZoneBatch(batch, zone, accessToken)
				if err != nil {
					fmt.Println("Error processing batch:", err)
					return
				}
				allZoneResults = append(allZoneResults, results)
			}
		}

		// Step 4: Write results to CSV at the end of 30 minutes
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
