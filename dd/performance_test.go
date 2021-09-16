package dd_test

/*
This example illustrates the performance 51Degrees device detection solution.
*/

import ( //	"runtime"
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/51Degrees/device-detection-go/dd"
)

// File to output the performance report
const reportFile = "performance_report.log"

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
		results, err := dd.NewResultsHash(
			manager,
			1,
			0)
		if err != nil {
			log.Fatal(err)
		}

		// Make sure results object is freed after function execution.
		defer func() {
			err = results.Free()
			if err != nil {
				panic(err)
			}
		}()

		// fmt.Println(ua)
		// Perform detection
		err = results.MatchUserAgent(ua)
		if err != nil {
			log.Fatal(err)
		}

		// Get the value in string
		value, _, err := results.GetValuesString(
			"IsMobile",
			100,
			",")
		if err != nil {
			log.Fatal(err)
		}

		// Update report
		if strings.Compare("True", value) == 0 {
			atomic.AddUint64(&rep.uaIsMobile, 1)
		}
	}

	// Complete and mark as done
	defer wg.Done()
}

// Count the number of User-Agents in a User-Agents file and update
// a report statistic
func countUAFromFiles(
	uaFilePath string,
	rep *report) {
	// Count the number of User Agents
	f, err := os.OpenFile(uaFilePath, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Count the number of UA.
	s := bufio.NewScanner(f)
	defer func() {
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	// Count the User-Agents
	for s.Scan() {
		rep.uaCount++
	}
}

// Run the performance test. Determine the number of records in a User-Agent
// file. Iterate through the User-Agent file and perform detection on each
// User-Agent. Record the processing time and update a report statistic.
func performDetections(
	manager *dd.ResourceManager,
	uaFilePath string,
	calibration bool,
	rep *report) {
	countUAFromFiles(uaFilePath, rep)
	// Loop through the User-Agent file
	file, err := os.OpenFile(uaFilePath, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure to close the file
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Actual processing
	scanner := bufio.NewScanner(file)
	defer func() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	// Create a wait group
	var wg sync.WaitGroup
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
	// Wait until all goroutines finish
	wg.Wait()
}

// Check a error returned from writing to a buffer
func checkWriteError(err error) {
	if err != nil {
		log.Fatal("ERROR: Failed to write to buffer.")
	}
}

// Print report to a report file
func printReport(caliR *report, actR *report) {
	// Check if a report file already exists
	if _, err := os.Stat(reportFile); err == nil || errors.Is(err, fs.ErrExist) {
		// If no 'force' option is specified then terminate.
		if !strings.EqualFold(os.Args[len(os.Args)-1], "force") {
			log.Fatalf("ERROR: A report file \"%s\" already exists.", reportFile)
		}
	}

	// Create a report file
	f, err := os.Create(reportFile)
	if err != nil {
		log.Fatalf("ERROR: Failed to create report file \"%s\".", reportFile)
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
}

// Run the performance example. Performs two phase: calibration and actual
// detection. Processing time of each phase is recorded to produce the actual
// processing time per detection
func run(
	manager *dd.ResourceManager,
	uaFilePath string) {
	// Calibration
	caliReport := report{0, 0, 0, 0}
	start := time.Now()
	performDetections(manager, uaFilePath, true, &caliReport)
	end := time.Now()
	caliTime := end.Sub(start)
	caliReport.processingTime = caliTime.Milliseconds()
	// Validation to make sure same number of UAs have been read and processed
	if caliReport.uaCount != caliReport.uaProcessed {
		log.Fatal("ERROR: Not all User-Agents have been processed.")
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
		log.Fatal("ERROR: Not all User-Agents have been processed.")
	}

	// Print the final performance report
	printReport(&caliReport, &actReport)
}

// Setup all configuration settings required for running this example.
// Run the example.
func runPerformanceExample(
	dataFilePath string,
	uaFilePath string,
	perf dd.PerformanceProfile) {
	// Create Resource Manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash()
	config.SetPerformanceProfile(perf)
	config.SetConcurrency(uint16(runtime.NumCPU()))
	config.SetUsePredictiveGraph(false)
	config.SetUsePredictiveGraph(true)
	config.SetUseUpperPrefixHeaders(false)
	config.SetUpdateMatchedUserAgent(false)
	err := dd.InitManagerFromFile(
		manager,
		config,
		"IsMobile",
		dataFilePath)
	if err != nil {
		panic(err)
	}

	// Make sure manager object will be freed after the function execution
	defer func() {
		err := manager.Free()
		if err != nil {
			panic(err)
		}
	}()

	// Run the performance tests
	run(manager, uaFilePath)
}

func Example_Performance() {
	// Data file path
	dataFilePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	// User-Agents file path
	uaFilePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.csv"

	runPerformanceExample(dataFilePath, uaFilePath, dd.Balanced)
	fmt.Printf("FINISHED")

	// The performance is output to a file 'performance_report.log' with content
	// similar as below:
	//   Average 0.01510 ms per User-Agent
	//   Total User-Agents: 20000
	//   IsMobile User-Agents: 14527
	//   Processed User-Agents: 20000
	//   Number of CPUs: 2

	// Output:
	// FINISHED
}
