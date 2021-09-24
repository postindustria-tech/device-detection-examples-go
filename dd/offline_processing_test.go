package dd_test

/*
This example illustrates how to perform simple device detections on given
User-Agent strings.
*/

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/51Degrees/device-detection-go/dd"
)

// function match performs a match on an input User-Agent string and determine
// if the device is a mobile device.
func processUserAgent(
	results *dd.ResultsHash,
	ua string) {
	// Perform detection
	err := results.MatchUserAgent(ua)
	if err != nil {
		log.Fatalf(
			"ERROR: \"%s\" on User-Agent \"%s\".\n", err, ua)
	}
}

func process(
	manager *dd.ResourceManager,
	uaFilePath string,
	outputFilePath string) {
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("ERROR: Failed to create file %s.\n", outputFilePath)
	}

	// Create results
	results := dd.NewResultsHash(manager, 1, 0)

	// Make sure results object is freed after function execution.
	defer results.Free()

	// Open the User-Agents file for processing
	uaFile, err := os.OpenFile(uaFilePath, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatalf("ERROR: Failed to open file \"%s\".\n", uaFilePath)
	}

	// Create a scanner to read User-Agents
	scanner := bufio.NewScanner(uaFile)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("ERROR: Failed during scanning file \"%s\".\n", uaFilePath)
		}
	}()

	w := io.Writer(outFile)
	available := results.AvailableProperties()

	// Print header to output file
	fmt.Fprintf(w, "\"User-Agent\",\"Drift\",\"Difference\",\"Iterations\"")
	for i := 0; i < len(available); i++ {
		fmt.Fprintf(w, ",\"%s\"", available[i])
	}
	fmt.Fprintf(w, "\n")

	for scanner.Scan() {
		processUserAgent(results, scanner.Text())
		// Output the matched  user agent string, drift, difference, iterations
		fmt.Fprintf(
			w,
			"\"%s\",%d,%d,%d",
			results.UserAgent(0),
			results.DriftByIndex(0),
			results.DifferenceByIndex(0),
			results.IterationsByIndex(0))
		// Get the values in string
		for i := 0; i < len(available); i++ {
			hasValues, err := results.HasValuesByIndex(i)
			if err != nil {
				log.Fatalln(err)
			}

			// Write empty value if one isn't available
			if !hasValues {
				fmt.Fprintf(w, ",\"\"")
			} else {
				value, err := results.ValuesString(
					"IsMobile",
					",")
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Fprintf(w, ",%s", value)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func runOfflineProcessing(perf dd.PerformanceProfile) string {
	// Initialise manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash(perf)
	filePath := getFilePath([]string{liteDataFile})
	uaFilePath := getFilePath([]string{uaFile})
	uaDir := filepath.Dir(uaFilePath)
	uaBase := filepath.Base(uaFilePath)
	outputFilePath := fmt.Sprintf("%s/%s.processed.csv", uaDir, uaBase)
	// Get base path
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed to get current directory.")
	}
	// Get relative output path for testing
	relOutputFilePath, err := filepath.Rel(basePath, outputFilePath)
	if err != nil {
		log.Fatalln("Failed to get relative output file path.")
	}

	config.SetUpdateMatchedUserAgent(true)
	err = dd.InitManagerFromFile(
		manager,
		*config,
		"IsMobile,BrowserName,DeviceType,PriceBand,ReleaseMonth,ReleaseYear",
		filePath)
	if err != nil {
		log.Fatalln(err)
	}

	// Make sure manager object will be freed after the function execution
	defer manager.Free()

	process(manager, uaFilePath, outputFilePath)
	return fmt.Sprintf("Output to \"%s\".\n", relOutputFilePath)
}

func Example_offline_processing() {
	performExample(dd.Default, runOfflineProcessing)
	// Output:
	// Output to "../device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.csv.processed.csv".
}
