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

package main

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
Average 0.00456 ms per Evidence Record
Total Evidence Records: 80000
IsMobile Evidence Records: 58076
Processed Evidence Records: 80000
Number of CPUs: 2
```
*/

import ( //	"runtime"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	dd_example "github.com/51Degrees/device-detection-examples-go/v4/dd"
	"gopkg.in/yaml.v3"

	"github.com/51Degrees/device-detection-go/v4/dd"
)

// File to output the performance report
const reportFile = "performance_report.log"

// Report struct for each performance run
type report struct {
	evidenceCount     uint64
	evidenceIsMobile  uint64
	evidenceProcessed uint64
	processingTime    int64
}

// Perform device detection on a Evidence Record
func matchEvidenceRecord( //TODO rewrite
	wg *sync.WaitGroup,
	manager *dd.ResourceManager,
	evidence *dd.Evidence,
	calibration bool,
	rep *report) {
	defer evidence.Free()
	// Increase the number of Evidence Record being processed
	atomic.AddUint64(&rep.evidenceProcessed, 1)
	if !calibration {
		results := dd.NewResultsHash(manager, uint32(evidence.Count()), 0)

		// Make sure results object is freed after function execution.
		defer results.Free()

		// Perform detection
		err := results.MatchEvidence(evidence)
		if err != nil {
			log.Fatal("ERROR: Failed to perform detection.")
		}

		// Get the value in string
		res, err := results.ValuesString(
			"IsMobile",
			",")
		if err != nil {
			log.Fatalln(err)
		}

		// Update report
		if strings.Compare("True", res) == 0 {
			atomic.AddUint64(&rep.evidenceIsMobile, 1)
		}
	}

	// Complete and mark as done
	defer wg.Done()
}

// Run the performance test. Determine the number of records in a Evidence
// file. Iterate through the Evidence file and perform detection on each
// Evidence. Record the processing time and update a report statistic.
func performDetections(
	manager *dd.ResourceManager,
	options dd_example.Options,
	calibration bool,
	rep *report) {
	// Create a wait group
	var wg sync.WaitGroup
	evidenceFilePath := dd_example.GetFilePathByPath(options.EvidenceFilePath)

	// rep.uaCount = dd_example.CountUAFromFiles(uaFilePath) TODO remove
	// rep.uaCount *= options.Iterations

	for i := 0; i < int(options.Iterations); i++ {
		// Loop through the Evidence file
		file, err := os.OpenFile(evidenceFilePath, os.O_RDONLY, 0444)
		if err != nil {
			log.Fatalf("ERROR: Failed to open file \"%s\".\n", evidenceFilePath)
		}
		defer func() {
			// Make sure the file is closed properly
			if err := file.Close(); err != nil {
				log.Fatalf("ERROR: Failed to close file \"%s\".\n", evidenceFilePath)
			}
		}()

		// Actual processing
		dec := yaml.NewDecoder(file)
		for {
			// Decode Evidence file by line
			var doc map[string]string
			if err := dec.Decode(&doc); err == io.EOF {
				break
			} else if err != nil {
				// Make sure there is no decoder error
				log.Fatalf("ERROR: Failed during decoding file \"%s\". %v\n", evidenceFilePath, err)
			}
			// Increase wait group
			wg.Add(1)
			rep.evidenceCount += 1

			// Prepare evidence for usage
			filteredEvidence := dd_example.ConvertEvidenceMap(doc)
			evidence := dd_example.ExtractEvidence(filteredEvidence)

			go matchEvidenceRecord(
				&wg,
				manager,
				evidence,
				calibration,
				rep)
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
func printReport(caliR *report, actR *report, logOutputPath string) string {
	// Get relative output path for testing
	var path string
	if filepath.IsAbs(logOutputPath) {
		path = logOutputPath
	} else {
		rootDir, e := os.Getwd()
		if e != nil {
			log.Fatalln("Failed to get current directory.")
		}
		path = filepath.Join(rootDir, logOutputPath)
	}
	path = filepath.Join(path, reportFile)

	// Create a report file
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("ERROR: Failed to create report file \"%s\".", path)
	}
	defer f.Close()

	// Create a writer
	w := bufio.NewWriter(f)

	// Make sure calibration and actual detections were performed on the same
	// number of Evidence Records.
	if actR.evidenceCount != caliR.evidenceCount {
		log.Fatal("ERROR: Calibration and actual detections were not" +
			"performed on the same number of Evidence Records.")
	}

	// Calculate actual performance
	avg := float64(actR.processingTime-caliR.processingTime) /
		float64(actR.evidenceCount)
	_, err = fmt.Fprintf(w, "Average %.5f ms per Evidence Record\n", avg)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "Total Evidence Records: %d\n", actR.evidenceCount)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "IsMobile Evidence Records: %d\n", actR.evidenceIsMobile)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "Processed Evidence Records: %d\n", actR.evidenceProcessed)
	checkWriteError(err)
	_, err = fmt.Fprintf(w, "Number of CPUs: %d\n", runtime.NumCPU())
	checkWriteError(err)
	w.Flush()
	return fmt.Sprintf("Output report to file \"%s\".\n", reportFile)
}

// Run the performance example. Performs two phase: calibration and actual
// detection. Processing time of each phase is recorded to produce the actual
// processing time per detection. Return output messages.
func run(
	manager *dd.ResourceManager,
	options dd_example.Options) string {
	// Calibration
	caliReport := report{0, 0, 0, 0}
	start := time.Now()
	performDetections(manager, options, true, &caliReport)
	end := time.Now()
	caliTime := end.Sub(start)
	caliReport.processingTime = caliTime.Milliseconds()
	// Validation to make sure same number of UAs have been read and processed
	if caliReport.evidenceCount != caliReport.evidenceProcessed {
		log.Fatalln("ERROR: Not all Evidence Records have been processed.")
	}

	// Action
	actReport := report{0, 0, 0, 0}
	start = time.Now()
	performDetections(manager, options, false, &actReport)
	end = time.Now()
	actTime := end.Sub(start)
	actReport.processingTime = actTime.Milliseconds()
	// Validation to make sure same number of UAs have been read and processed
	if actReport.evidenceCount != actReport.evidenceProcessed {
		log.Fatalln("ERROR: Not all Evidence Records have been processed.")
	}

	// Print the final performance report
	return printReport(&caliReport, &actReport, options.LogOutputPath)
}

// Setup all configuration settings required for running this example.
// Run the example.
func runPerformance(perf dd.PerformanceProfile, options dd_example.Options) string {
	dataFilePath := dd_example.GetFilePathByPath(options.DataFilePath)

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
	return run(manager, options)
}

func main() {
	dd_example.PerformExampleOptions(dd.InMemory, runPerformance)
	// The performance is output to a file 'performance_report.log' with content
	// similar as below:
	//   Average 0.01510 ms per Evidence Record
	//   Total Evidence Records: 20000
	//   IsMobile Evidence Records: 14527
	//   Processed Evidence Records: 20000
	//   Number of CPUs: 2

	// Output:
	// Output report to file "performance_report.log".
}
