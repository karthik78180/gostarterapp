package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Entry struct {
	SnNo int    `json:"sn no"`
	Host string `json:"host"`
	Port string `json:"port"`
}

type Result struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Zone         string `json:"zone"`
	Status       string `json:"status"`
	ResponseBody string `json:"response_body"`
	Error        string `json:"error,omitempty"`
}

var zones = []string{"zone1", "zone2", "zone3", "zone4"}
var accessToken string

func refreshToken() string {
	// Logic to refresh the access token
	// You can call the API here to get the new token and return it
	// Placeholder token for demonstration
	return "newAccessToken"
}

func callAPI(zone, host, port string) Result {
	url := fmt.Sprintf("http://api.com/connectivitycheck?zone=%s&host=%s&port=%s", zone, host, port)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Result{Host: host, Port: port, Zone: zone, Status: "failure", Error: err.Error()}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Result{Host: host, Port: port, Zone: zone, Status: "failure", Error: err.Error()}
	}

	if resp.StatusCode == http.StatusOK {
		return Result{Host: host, Port: port, Zone: zone, Status: "success", ResponseBody: string(body)}
	}
	return Result{Host: host, Port: port, Zone: zone, Status: "failure", ResponseBody: string(body), Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
}

func processBatch(entries []Entry, zone string, resultsChan chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	var wgBatch sync.WaitGroup
	for _, entry := range entries {
		wgBatch.Add(1)
		go func(e Entry) {
			defer wgBatch.Done()
			result := callAPI(zone, e.Host, e.Port)
			resultsChan <- result
		}(entry)
	}
	wgBatch.Wait()
}

func writeResultsToCSV(results [][]Result, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	headers := []string{"host", "port"}
	for _, zone := range zones {
		headers = append(headers, fmt.Sprintf("%s status", zone), fmt.Sprintf("%s response body", zone), fmt.Sprintf("%s error", zone))
	}
	err = writer.Write(headers)
	if err != nil {
		return err
	}

	// Write the rows
	for i := 0; i < len(results[0]); i++ {
		row := []string{results[0][i].Host, results[0][i].Port}
		for _, zoneResults := range results {
			row = append(row, zoneResults[i].Status, zoneResults[i].ResponseBody, zoneResults[i].Error)
		}
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Reading the input JSON file
	inputFile := "input.json"
	jsonData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %s\n", err)
		return
	}

	var entries []Entry
	err = json.Unmarshal(jsonData, &entries)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s\n", err)
		return
	}

	// Output file that will be written every 30 minutes
	outputFile := fmt.Sprintf("output_%s.csv", strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "_"))

	// Prepare results channel and waitgroup
	for {
		accessToken = refreshToken()
		var zoneResults [][]Result

		for _, zone := range zones {
			resultsChan := make(chan Result, len(entries))
			var wg sync.WaitGroup

			// Process entries in batches of 10 for the current zone
			for i := 0; i < len(entries); i += 10 {
				end := i + 10
				if end > len(entries) {
					end = len(entries)
				}
				wg.Add(1)
				go processBatch(entries[i:end], zone, resultsChan, &wg)

				// Wait for 30 seconds between batches
				if end != len(entries) {
					time.Sleep(30 * time.Second)
				}
			}

			// Wait for all batches to finish
			wg.Wait()
			close(resultsChan)

			// Collect results for the current zone
			var results []Result
			for result := range resultsChan {
				results = append(results, result)
			}
			zoneResults = append(zoneResults, results)
		}

		// Write results to CSV after processing all zones
		err = writeResultsToCSV(zoneResults, outputFile)
		if err != nil {
			fmt.Printf("Error writing CSV file: %s\n", err)
			return
		}

		fmt.Println("Results written to", outputFile)

		// Wait for 30 minutes before refreshing the token and processing the next batch
		time.Sleep(30 * time.Minute)
	}
}
