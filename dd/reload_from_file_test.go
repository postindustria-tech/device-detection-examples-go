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
Illustrates how dataset can be reloaded while detections are performed.
*/

import (
	"bufio"
	"hash/fnv"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/51Degrees/device-detection-go/v4/dd"
)

// Number of iterations to perform over the User-Agents.
const fIterationCount = 4

// Report struct for reload from file rn
type freport struct {
	mu          sync.Mutex // Mutex
	uaCount     uint64
	hashCodes   [fIterationCount]uint32
	uaProcessed uint64
}

// updateHashCode updates the hash code with the input code ad the index
// specified. The update use XOR operation. This function is thread safe to
// make sure multiple threads can update the hash code correctly
func (rep *freport) updateHashCode(code uint32, i uint32) {
	rep.mu.Lock()
	rep.hashCodes[i] ^= code
	rep.mu.Unlock()
}

// generateHash generate 32bit hash code for an input string
func generateHash(str string) uint32 {
	h := fnv.New32()
	h.Write([]byte(str))
	return h.Sum32()
}

func executeTest(
	wg *sync.WaitGroup,
	manager *dd.ResourceManager,
	ua string,
	rep *freport,
	iteration uint32) {
	// Create results
	results := dd.NewResultsHash(manager, 1, 0)

	// Make sure results object is freed after function execution.
	defer results.Free()

	// Perform detection
	results.MatchUserAgent(ua)

	// Loop through all properties
	for _, property := range results.AvailableProperties() {
		// Get the value in string
		value, err := results.ValuesString(
			property,
			",")
		if err != nil {
			log.Fatalln(err)
		}
		rep.updateHashCode(generateHash(value), iteration)
	}

	// Increase the number of User-Agents processed
	atomic.AddUint64(&rep.uaProcessed, 1)

	// Complete and mark as done
	defer wg.Done()
}

// performDetectionInterations iterates through the User-Agents file and perform
// detection on each User-Agent. Results of each detection will be hashed and
// combine for each iteration. At the end all itertions should have the same
// hash value. If the hash values are different, it indicates that User-Agents
// might have not been processed correctly in some iterations.
func performDetectionIterations(
	manager *dd.ResourceManager,
	uaFilePath string,
	wg *sync.WaitGroup,
	rep *freport) {
	for i := 0; i < fIterationCount; i++ {
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
			go executeTest(
				wg,
				manager,
				scanner.Text(),
				rep,
				uint32(i))
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
	wg.Done()
}

func runReloadFromFileSub(
	manager *dd.ResourceManager,
	uaFilePath string) string {
	reloads := 0
	reloadFails := 0
	// Create a wait group for iteration function
	var wg sync.WaitGroup

	// Count the number of User-Agents to be processed
	var rep freport
	rep.uaCount = countUAFromFiles(uaFilePath)
	rep.uaCount *= fIterationCount

	// Perform detections
	wg.Add(1)
	go performDetectionIterations(manager, uaFilePath, &wg, &rep)

	// Perform reload from file until all User-Agents have been processed
	for rep.uaProcessed < rep.uaCount {
		err := manager.ReloadFromOriginalFile()
		if err == nil {
			// Failed to reload the original file
			reloads++
		} else {
			reloadFails++
		}
		// Sleep 1 second between reload
		time.Sleep(1000000000) // This is in nano seconds
	}

	// Wait until all goroutines finish
	wg.Wait()

	// Construct report
	log.Printf("Reloaded '%d' times.\n", reloads)
	log.Printf("Failed to reload '%d' times.\n", reloadFails)
	var initHashCode uint32
	for i := 0; i < fIterationCount; i++ {
		if i == 0 {
			initHashCode = rep.hashCodes[i]
		} else if initHashCode != rep.hashCodes[i] {
			log.Fatalf("Hash codes do not match. Initial hash code is '%d', "+
				"but iteration '%d' has hash code '%d'. This indicates not "+
				"all User-Agents have been processed correctly for each "+
				"iteration.", initHashCode, rep.hashCodes[i], i)
		}
		log.Printf("Hashcode '%d' for iteration '%d'.\n",
			rep.hashCodes[i], i)
	}
	return "Program execution complete."
}

func runReloadFromFile(perf dd.PerformanceProfile) string {
	dataFilePath := getFilePath([]string{liteDataFile})
	uaFilePath := getFilePath([]string{uaFile})

	// Create Resource Manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash(dd.InMemory)
	config.SetConcurrency(uint16(runtime.NumCPU()))
	config.SetUseUpperPrefixHeaders(false)
	config.SetUpdateMatchedUserAgent(false)
	err := dd.InitManagerFromFile(
		manager,
		*config,
		"IsMobile,BrowserName,DeviceType",
		dataFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// Make sure manager object will be freed after the function execution
	defer manager.Free()

	// Run the performance tests
	report := runReloadFromFileSub(manager, uaFilePath)
	return report
}

func Example_reload_from_file() {
	performExample(dd.Default, runReloadFromFile)
	// The output log of this example is in for the following format:
	//
	// 2021/11/10 11:42:05 Reloaded '2' times.
	// 2021/11/10 11:42:05 Failed to reload '0' times.
	// 2021/11/10 11:42:05 Hashcode '4217895257' for iteration '0'.
	// 2021/11/10 11:42:05 Hashcode '4217895257' for iteration '1'.
	// 2021/11/10 11:42:05 Hashcode '4217895257' for iteration '2'.
	// 2021/11/10 11:42:05 Hashcode '4217895257' for iteration '3'.

	// Output:
	// Program execution complete.
}
