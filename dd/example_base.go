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

package dd_example

/*
This example illustrates how to perform simple device detections on given
User-Agent strings.
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/51Degrees/device-detection-go/v4/dd"
)

// Constants
const LiteDataFile = "51Degrees-LiteV4.1.hash"
const EnterpriseDataFile = "Enterprise-HashV41.hash"
const UaFile = "20000 User Agents.csv"

// Type take a performance profile, run the code and get the return output
type ExampleFunc func(p dd.PerformanceProfile) string

// Returns a full path to a file to be used for examples
func GetFilePath(names []string) string {
	filePath, err := dd.GetFilePath(
		"..",
		names,
	)
	if err != nil {
		log.Fatalf("Could not find any file that matches any of \"%s\".\n",
			strings.Join(names, ", "))
	}
	return filePath
}

// isFlagOn checks if certain flag is enabled in the test input
// args.
func isFlagOn(value string) bool {
	for _, arg := range os.Args {
		if strings.EqualFold(value, arg) {
			return true
		}
	}
	return false
}

// Count the number of User-Agents in a User-Agents file and return the number
// of user agents found.
func CountUAFromFiles(
	uaFilePath string) uint64 {
	var count uint64 = 0
	// Count the number of User Agents
	f, err := os.OpenFile(uaFilePath, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatalf("ERROR: Failed to open file \"%s\".\n", uaFilePath)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("ERROR: Failed to close file \"%s\".\n", uaFilePath)
		}
	}()

	// Count the number of UA.
	s := bufio.NewScanner(f)
	defer func() {
		if err := s.Err(); err != nil {
			log.Fatalf("ERROR: Error during scanning file\"%s\".\n", uaFilePath)
		}
	}()

	// Count the User-Agents
	for s.Scan() {
		count++
	}
	return count
}

// This is a wrapper function which execute a function that contains
// example code with an input performance profile or all performance
// profiles if performed under CI.
func PerformExample(perf dd.PerformanceProfile, eFunc ExampleFunc) {
	perfs := []dd.PerformanceProfile{perf}
	// If running under ci, use all performance profiles
	if isFlagOn("ci") {
		perfs = []dd.PerformanceProfile{
			dd.Default,
			dd.LowMemory,
			dd.Balanced,
			dd.BalancedTemp,
			dd.HighPerformance,
			dd.InMemory,
		}
	}
	// Execute the example function with all performance profiles
	for i, p := range perfs {
		output := eFunc(p)
		// This is to support example Output verification
		// so only print once.
		if i == 0 {
			fmt.Print(output)
		}
	}
}
