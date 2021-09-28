/* *********************************************************************
 * This Original Work is copyright of 51 Degrees Mobile Experts Limited.
 * Copyright 2019 51 Degrees Mobile Experts Limited, 5 Charlotte Close,
 * Caversham, Reading, Berkshire, United Kingdom RG4 7BY.
 *
 * This Original Work is licensed under the European Union Public Licence (EUPL)
 * v.1.2 and is subject to its terms as set out below.
 *
 * If a copy of the EUPL was not distributed with this file, You can obtain
 * one at https://opensource.org/licenses/EUPL-1.2.
 *
 * The 'Compatible Licences' set out in the Appendix to the EUPL (as may be
 * amended by the European Commission) shall be deemed incompatible for
 * the purposes of the Work and the provisions of the compatibility
 * clause in Article 5 of the EUPL shall not apply.
 *
 * If using the Work as, or as part of, a network application, by
 * including the attribution notice(s) required under Article 5 of the EUPL
 * in the end user terms of the application under an appropriate heading,
 * such notice(s) shall fulfill the requirements of that article.
 * ********************************************************************* */

package dd_test

/*
This example illustrates how to process a list of User-Agents from file and
output detection metrics and properties of each User-Agent to another file for
further evaluation.

Expected output is as described at the "// Output:..." section locate at the
bottom of this example.

To run this example, perform the following command:
```
go test -run Example_offline_processing
```

This example will output to a file located at
"../device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.processed.csv".
This contains the detection metrics User-Agent,Drift,Difference,Iterations and
available properties for each User-Agent.

*/

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	uaBase := strings.TrimSuffix(filepath.Base(uaFilePath), filepath.Ext(uaFilePath))
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
	// Output to "../device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.processed.csv".
}
