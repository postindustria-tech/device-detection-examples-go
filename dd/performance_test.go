package dd_test

/*
This example illustrates the performance 51Degrees device detection solution.
*/

import ( //	"runtime"
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/51Degrees/device-detection-go/dd"
)

const progressMarks = 40

type report struct {
	uaCount        uint64
	uaIsMobile     uint64
	uaProcessed    uint64
	processingTime int64
}

// TODO: This does not work nicely in multi threads environment
// Don't print load bar for now as it needs to be done by one
// thread which cannot be determined how often the work is allocated for
// that thread within golang concurrency model.
/*
func printLoadBar(r *report) {
	processed := r.uaProcessed
	progress := r.uaCount / progressMarks
	full := processed / progress
	empty := r.uaCount/progress - full
	fmt.Printf("\r[")
	for i := uint64(0); i < full; i++ {
		fmt.Printf("=")
	}
	for i := uint64(0); i < empty; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("]")
}
*/

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

func runTests(
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

func printReport(r *report) {
	fmt.Printf("\n")
	fmt.Printf("Total User-Agents: %d\n", r.uaCount)
	fmt.Printf("IsMobile User-Agents: %d\n", r.uaIsMobile)
	fmt.Printf("Processed User-Agents: %d\n", r.uaProcessed)
	fmt.Printf("Processing time: %d ms\n", r.processingTime)
}

func run(
	manager *dd.ResourceManager,
	uaFilePath string) {
	// Calibration
	fmt.Println("\nCalibrating...")
	caliReport := report{0, 0, 0, 0}
	start := time.Now()
	runTests(manager, uaFilePath, true, &caliReport)
	end := time.Now()
	caliTime := end.Sub(start)
	caliReport.processingTime = caliTime.Milliseconds()
	printReport(&caliReport)
	// Validation to make sure same number of UAs have been read and processed
	if caliReport.uaCount != caliReport.uaProcessed {
		log.Fatal("ERROR: Not all User-Agents have been processed.")
	}

	// Action
	fmt.Println("\nRunning performance tests...")
	actReport := report{0, 0, 0, 0}
	start = time.Now()
	runTests(manager, uaFilePath, false, &actReport)
	end = time.Now()
	actTime := end.Sub(start)
	actReport.processingTime = actTime.Milliseconds()
	printReport(&actReport)
	// Validation to make sure same number of UAs have been read and processed
	if actReport.uaCount != actReport.uaProcessed {
		log.Fatal("ERROR: Not all User-Agents have been processed.")
	}

	// Make sure calibration and actual detections were performed on the same
	// number of user agents.
	if actReport.uaCount != caliReport.uaCount {
		log.Fatal("ERROR: Calibration and actual detections were not" +
			"performed on the same number of User-Agents")
	}

	// Calculate actual performance
	avg := float64(
		actTime.Milliseconds()-caliTime.Milliseconds()) /
		float64(actReport.uaCount)
	fmt.Printf("Average %.4f ms per User-Agent\n", avg)
}

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
	// Output:
	// .*
}
