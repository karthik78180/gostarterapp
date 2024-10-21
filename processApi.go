package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Entry struct {
	SnNo int    `json:"sn no"`
	Host string `json:"host"`
	Port string `json:"port"`
}

type ApiResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Result struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Zone         string `json:"zone"`
	Status       string `json:"status"`
	ResponseBody string `json:"response_body"`
	Error        string `json:"error,omitempty"`
}

// ProcessZoneBatch processes a single batch for a specific zone and access token
func ProcessZoneBatch(entries []Entry, zone, accessToken string) ([]Result, error) {
	resultsChan := make(chan Result, len(entries))
	var wg sync.WaitGroup

	wg.Add(1)
	go processBatch(entries, zone, resultsChan, &wg, accessToken)

	wg.Wait()
	close(resultsChan)

	// Collect results
	var results []Result
	for result := range resultsChan {
		results = append(results, result)
	}

	return results, nil
}

// processBatch processes a batch of entries and makes API calls in parallel
func processBatch(entries []Entry, zone string, resultsChan chan<- Result, wg *sync.WaitGroup, accessToken string) {
	defer wg.Done()

	for _, entry := range entries {
		wg.Add(1)
		go func(e Entry) {
			defer wg.Done()

			// Call the API and handle the response
			res, resBody, err := callApi(e.Host, e.Port, zone, accessToken)
			if err != nil {
				resultsChan <- Result{Host: e.Host, Port: e.Port, Zone: zone, Status: "failure", ResponseBody: resBody, Error: err.Error()}
				return
			}

			// Successful API response
			if res.Status == "success" && strings.Contains(res.Message, "istio-proxy : Connectivity result from container istio-proxy is successful.") {
				resultsChan <- Result{Host: e.Host, Port: e.Port, Zone: zone, Status: "success", ResponseBody: resBody}
			} else {
				resultsChan <- Result{Host: e.Host, Port: e.Port, Zone: zone, Status: "failure", ResponseBody: resBody}
			}
		}(entry)
	}
}

// callApi simulates an API call that returns an ApiResponse and an error
func callApi(host, port, zone, accessToken string) (ApiResponse, string, error) {
	url := fmt.Sprintf("http://api.com/connectivitycheck?zone=%s&host=%s&port=%s", zone, host, port)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ApiResponse{}, "", err
	}
	defer resp.Body.Close()

	var apiResponse ApiResponse
	resBody := "Response body here" // You would read the response body here
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return ApiResponse{}, resBody, err
	}

	return apiResponse, resBody, nil
}

// writeResultsToCSV writes the results to a CSV file
func writeResultsToCSV(allZoneResults [][]Result, outputFile string) error {
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
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write the rows
	for i := 0; i < len(allZoneResults[0]); i++ {
		row := []string{allZoneResults[0][i].Host, allZoneResults[0][i].Port}
		for _, zoneResults := range allZoneResults {
			row = append(row, zoneResults[i].Status, zoneResults[i].ResponseBody, zoneResults[i].Error)
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Dummy function to load entries
func loadEntries() []Entry {
	return []Entry{
		{SnNo: 1, Host: "example.com", Port: "80"},
		{SnNo: 2, Host: "example.com", Port: "443"},
		{SnNo: 3, Host: "example2.com", Port: "443"},
	}
}
