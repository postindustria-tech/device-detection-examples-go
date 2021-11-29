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
This example illustrates the performance of 51Degrees device detection solution.

Expected output is as described at the "// Output:..." section locate at the
bottom of this example.

To run this example, perform the following command:
```
go test -run Example_performance
```

This example will output a report to ./performance_report.log. The report
content is in the below format:
```
Average 0.00456 ms per User-Agent
Total User-Agents: 80000
IsMobile User-Agents: 58076
Processed User-Agents: 80000
Number of CPUs: 2
```
*/

import ( //	"runtime"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/51Degrees/device-detection-go/v4/dd"
)

// File to output the performance report
const reportFile = "performance_report.log"

// Number of iterations to perform over the User-Agents.
// The higher the number is the more accurate the report.
const iterationCount = 4

// Report struct for each performance run
type report struct {
	uaCount        uint64
	uaIsMobile     uint64
	uaProcessed    uint64
	processingTime int64
}

// Perform device detection on a User-Agent
func matchUserAgent(
	wg *sync.WaitGroup,
	manager *dd.ResourceManager,
	ua string,
	calibration bool,
	rep *report) {
	// Increase the number of User-Agents being processed
	atomic.AddUint64(&rep.uaProcessed, 1)
	if !calibration {
		// Create results
		results := dd.NewResultsHash(manager, 1, 0)

		// Make sure results object is freed after function execution.
		defer results.Free()

		// fmt.Println(ua)
		// Perform detection
		results.MatchUserAgent(ua)

		// Get the value in string
		value, err := results.ValuesString(
			"IsMobile",
			",")
		if err != nil {
			log.Fatalln(err)
		}

		// Update report
		if strings.Compare("True", value) == 0 {
			atomic.AddUint64(&rep.uaIsMobile, 1)
		}
	}

	// Complete and mark as done
	defer wg.Done()
}

// Run the performance test. Determine the number of records in a User-Agent
// file. Iterate through the User-Agent file and perform detection on each
// User-Agent. Record the processing time and update a report statistic.
func performDetections(
	manager *dd.ResourceManager,
	uaFilePath string,
	calibration bool,
	rep *report) {
	// Create a wait group
	var wg sync.WaitGroup

	rep.uaCount = countUAFromFiles(uaFilePath)
	rep.uaCount *= iterationCount

	for i := 0; i < iterationCount; i++ {
		// Loop through the User-Agent file
		file, err := os.OpenFile(uaFilePath, os.O_RDONLY, 0444)
		if err != nil {
			log.Fatalf("ERROR: Failed to open file \"%s\".\n", uaFilePath)
		}

		// Actual processing
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			// Increase wait group
			wg.Add(1)
			go matchUserAgent(
				&wg,
				manager,
				scanner.Text(),
				calibration,
				rep)
		}

		// Make sure there is no scanner error
		if err := scanner.Err(); err != nil {
			log.Fatalf("ERROR: Error during scanning file \"%s\".\n", uaFilePath)
		}

		// Make sure the file is closed properly
		if err := file.Close(); err != nil {
			log.Fatalf("ERROR: Failed to close file \"%s\".\n", uaFilePath)
		}
	}
	// Wait until all goroutines finish
	wg.Wait()
}

// Check a error returned from writing to a buffer
func checkWriteError(err error) {
	if err != nil {
		log.Fatalln("ERROR: Failed to write to buffer.")
	}
}

// Print report to a report file and return output message.
func printReport(caliR *report, actR *report) string {
	// Get base path
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed to get current directory.")
	}
	reportFilePath := basePath + "/" + reportFile
	// Get relative output path for testing
	relReportFilePath, err := filepath.Rel(basePath, reportFilePath)
	if err != nil {
		log.Fatalln("Failed to get relative output file path.")
	}

	// Create a report file
	f, err := os.Create(reportFilePath)
	if err != nil {
		log.Fatalf("ERROR: Failed to create report file \"%s\".", reportFilePath)
	}
	defer f.Close()

	// Create a writer
	w := bufio.NewWriter(f)

	// Make sure calibration and actual detections were performed on the same
	// number of user agents.
	if actR.uaCount != caliR.uaCount {
		log.Fatal("ERROR: Calibration and actual detections were not" +
			"performed on the same number of User-Agents.")
	}

	// Calculate actual performance
	avg := float64(actR.processingTime-caliR.processingTime) /
		float64(actR.uaCount)
	_, err = fmt.Fprintf(w, "Average %.5f ms per User-Agent\n", avg)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "Total User-Agents: %d\n", actR.uaCount)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "IsMobile User-Agents: %d\n", actR.uaIsMobile)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "Processed User-Agents: %d\n", actR.uaProcessed)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "Number of CPUs: %d\n", runtime.NumCPU())
	checkWriteError(err)
	w.Flush()
	return fmt.Sprintf("Output report to file \"%s\".\n", relReportFilePath)
}

// Run the performance example. Performs two phase: calibration and actual
// detection. Processing time of each phase is recorded to produce the actual
// processing time per detection. Return output messages.
func run(
	manager *dd.ResourceManager,
	uaFilePath string) string {
	// Calibration
	caliReport := report{0, 0, 0, 0}
	start := time.Now()
	performDetections(manager, uaFilePath, true, &caliReport)
	end := time.Now()
	caliTime := end.Sub(start)
	caliReport.processingTime = caliTime.Milliseconds()
	// Validation to make sure same number of UAs have been read and processed
	if caliReport.uaCount != caliReport.uaProcessed {
		log.Fatalln("ERROR: Not all User-Agents have been processed.")
	}

	// Action
	actReport := report{0, 0, 0, 0}
	start = time.Now()
	performDetections(manager, uaFilePath, false, &actReport)
	end = time.Now()
	actTime := end.Sub(start)
	actReport.processingTime = actTime.Milliseconds()
	// Validation to make sure same number of UAs have been read and processed
	if actReport.uaCount != actReport.uaProcessed {
		log.Fatalln("ERROR: Not all User-Agents have been processed.")
	}

	// Print the final performance report
	return printReport(&caliReport, &actReport)
}

// Setup all configuration settings required for running this example.
// Run the example.
func runPerformance(perf dd.PerformanceProfile) string {
	dataFilePath := getFilePath([]string{liteDataFile})
	uaFilePath := getFilePath([]string{uaFile})

	// Create Resource Manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash(dd.InMemory)
	config.SetConcurrency(uint16(runtime.NumCPU()))
	config.SetUsePredictiveGraph(false)
	config.SetUsePerformanceGraph(true)
	config.SetUseUpperPrefixHeaders(false)
	config.SetUpdateMatchedUserAgent(false)
	err := dd.InitManagerFromFile(
		manager,
		*config,
		"IsMobile",
		dataFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// Make sure manager object will be freed after the function execution
	defer manager.Free()

	// Run the performance tests
	return run(manager, uaFilePath)
}

func Example_performance() {
	performExample(dd.InMemory, runPerformance)
	// The performance is output to a file 'performance_report.log' with content
	// similar as below:
	//   Average 0.01510 ms per User-Agent
	//   Total User-Agents: 20000
	//   IsMobile User-Agents: 14527
	//   Processed User-Agents: 20000
	//   Number of CPUs: 2

	// Output:
	// Output report to file "performance_report.log".
}
