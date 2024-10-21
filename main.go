package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Record struct {
	SnNo int    `json:"sn no"`
	Host string `json:"host"`
	Port string `json:"port"`
}

func readCSV(filePath string) ([]Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Check if 'host' and 'port' headers are present
	hostIndex, portIndex := -1, -1
	for i, header := range headers {
		if header == "host" {
			hostIndex = i
		} else if header == "port" {
			portIndex = i
		}
	}

	if hostIndex == -1 || portIndex == -1 {
		return nil, fmt.Errorf("CSV missing required headers: 'host' or 'port'")
	}

	var records []Record
	snCounter := 1

	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		host := line[hostIndex]
		ports := strings.Split(line[portIndex], "|") // split the ports by '|'

		for _, port := range ports {
			records = append(records, Record{
				SnNo: snCounter,
				Host: host,
				Port: port,
			})
			snCounter++
		}
	}

	return records, nil
}

func writeResultsToFile(records []Record, outputFilePath string) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	inputFilePath := flag.String("input", "resources/input.csv", "Path to input CSV file")
	outputFilePath := flag.String("output", "resources/output.json", "Path to output JSON file")
	flag.Parse()

	records, err := readCSV(*inputFilePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	err = writeResultsToFile(records, *outputFilePath)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}

	fmt.Println("CSV processing completed. Results written to", *outputFilePath)
}
