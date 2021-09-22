package dd_test

/*
This example illustrates how to perform simple device detections on given
User-Agent strings.
*/

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/51Degrees/device-detection-go/dd"
)

// Constants
const liteDataFile = "51Degrees-LiteV4.1.hash"
const enterpriseDataFile = "Enterprise-HashV41.hash"
const uaFile = "20000 User Agents.csv"

// Type take a performance profile, run the code and get the return output
type ExampleFunc func(p dd.PerformanceProfile) string

// Returns a full path to a file to be used for examples
func getFilePath(names []string) string {
	filePath, err := dd.GetFilePath(
		"../device-detection-go",
		names,
	)
	if err != nil && err != fs.ErrExist {
		fileList := ""
		for _, file := range names {
			if fileList != "" {
				fileList += ", "
			}
			fileList += file
		}
		log.Fatalf("Could not find any file that matches any of \"%s\".\n",
			fileList)
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

// This is a wrapper function which execute a function that contains
// example code with an input performance profile or all performance
// profiles if performed under CI.
func performExample(perf dd.PerformanceProfile, eFunc ExampleFunc) {
	perfs := []dd.PerformanceProfile{perf}
	// If running under ci, use all performance profiles
	if isFlagOn("ci") {
		perfs = []dd.PerformanceProfile{
			dd.Default,
			dd.LowMemory,
			dd.Balanced,
			// dd.BalancedTemp, // TODO: Enable once fixed
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
