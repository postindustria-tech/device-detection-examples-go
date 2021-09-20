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

	"github.com/51Degrees/device-detection-go/dd"
)

// function match performs a match on an input User-Agent string and determine
// if the device is a mobile device.
func processUserAgent(
	results dd.ResultsHash,
	ua string) {
	// Perform detection
	err := results.MatchUserAgent(ua)
	if err != nil {
		log.Fatalf(
			"ERROR: Failed to perform detection on User-Agent \"%s\".\n", ua)
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
	results, err := dd.NewResultsHash(
		manager,
		1,
		0)
	if err != nil {
		log.Fatalln("ERROR: Failed to create new results hash object.")
	}

	// Make sure results object is freed after function execution.
	defer func() {
		err = results.Free()
		if err != nil {
			log.Fatalln("ERROR: Ftailed to free results hash object.")
		}
	}()

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
	available, err := results.AvailableProperties()
	if err != nil {
		log.Fatalln("ERROR: Failed to obtain available properties.")
	}

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
			hasValues, err := results.HasValues(i)
			if err != nil {
				log.Fatal("ERROR: Failed to get 'HasValues'.")
			}

			// Write empty value if one isn't available
			if !hasValues {
				fmt.Fprintf(w, ",\"\"")
			} else {
				value, _, err := results.ValuesString(
					"IsMobile",
					100,
					",")
				if err != nil {
					log.Fatal("ERROR: Failed to get Values string.")
				}

				fmt.Fprintf(w, ",%s", value)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func Example_offline_processing() {
	// Initialise manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash(dd.Balanced)
	filePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	uaFilePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.csv"
	outputFilePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.processed.csv"
	config.SetUpdateMatchedUserAgent(true)

	err := dd.InitManagerFromFile(
		manager,
		*config,
		"IsMobile,BrowserName,DeviceType,PriceBand,ReleaseMonth,ReleaseYear",
		filePath)
	if err != nil {
		log.Fatalln("ERROR: Failed to initialize resource manager.")
	}

	// Make sure manager object will be freed after the function execution
	defer func() {
		err := manager.Free()
		if err != nil {
			log.Fatalln("ERROR: Failed to free resource manager.")
		}
	}()

	process(manager, uaFilePath, outputFilePath)
	fmt.Printf("FINISHED")

	// Output:
	// FINISHED
}
