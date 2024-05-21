package main

import (
	"hash/fnv"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	dd_example "github.com/51Degrees/device-detection-examples-go/v4/dd"

	"github.com/51Degrees/device-detection-examples-go/v4/onpremise/common"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"gopkg.in/yaml.v3"
)

// Number of iterations to perform over the Evidence Records.
const fIterationCount = 4

// Report struct for reload from file rn
type freport struct {
	mu                sync.Mutex // Mutex
	evidenceCount     uint64
	hashCodes         [fIterationCount]uint32
	evidenceProcessed uint64
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
	pl *onpremise.Engine,
	wg *sync.WaitGroup,
	evidence []onpremise.Evidence,
	rep *freport,
	iteration uint32) {
	results, err := pl.Process(evidence)
	if err != nil {
		log.Fatalln(err)
	}
	defer results.Free()

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

	// Increase the number of Evidence Records processed
	atomic.AddUint64(&rep.evidenceProcessed, 1)

	// Complete and mark as done
	defer wg.Done()
}

// performDetectionInterations iterates through the Evidence Records file and perform
// detection on each evidence. Results of each detection will be hashed and
// combine for each iteration. At the end all itertions should have the same
// hash value. If the hash values are different, it indicates that Evidence Records
// might have not been processed correctly in some iterations.
func performDetectionIterations(
	pl *onpremise.Engine,
	evidenceFilePath string,
	wg *sync.WaitGroup,
	rep *freport) {
	for i := 0; i < fIterationCount; i++ {
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
				log.Fatalf("ERROR: Error during decoding file \"%s\". %v\n", evidenceFilePath, err)
			}
			// Increase wait group
			wg.Add(1)

			// Prepare evidence for usage
			evidence := common.ConverToEvidence(doc)

			go executeTest(
				pl,
				wg,
				evidence,
				rep,
				uint32(i))
		}
	}
	wg.Done()
}

func runReloadFromFileSub(
	pl *onpremise.Engine,
	evidenceFilePath string, dataFilePath string) string {
	reloads := 0
	reloadFails := 0
	// Create a wait group for iteration function
	var wg sync.WaitGroup

	// Count the number of Evidence Records to be processed
	var rep freport
	rep.evidenceCount = dd_example.CountEvidenceFromFiles(evidenceFilePath)
	rep.evidenceCount *= fIterationCount

	// Perform detections
	wg.Add(1)
	go performDetectionIterations(pl, evidenceFilePath, &wg, &rep)

	// Perform reload from file until all Evidence Records have been processed
	for rep.evidenceProcessed < rep.evidenceCount {
		currentTime := time.Now().Local()
		err := os.Chtimes(dataFilePath, currentTime, currentTime)
		if err != nil {
			reloadFails++
		} else {
			reloads++
		}
		// Sleep 3 seconds between reload
		time.Sleep(3 * time.Second)
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
				"all Evidence Records have been processed correctly for each "+
				"iteration.", initHashCode, rep.hashCodes[i], i)
		}
		log.Printf("Hashcode '%d' for iteration '%d'.\n",
			rep.hashCodes[i], i)
	}
	return "Program execution complete."
}

func main() {
	//common.LoadEnvFile()

	common.RunExample(
		func(params common.ExampleParams) error {
			//... Example code
			//Create config
			config := dd.NewConfigHash(dd.InMemory)
			//config.SetConcurrency(20)

			//Create on-premise engine
			pl, err := onpremise.New(
				config,
				//Provide your own file URL
				//it can be compressed as gz or raw, engine will handle it
				// Passing a custom update URL
				// Path to your data file

				onpremise.WithDataFile(params.DataFile),
				// For automatic updates to work you will need to provide a license key.
				// A license key can be obtained with a subscription from https://51degrees.com/pricing
				onpremise.WithLicenseKey(params.LicenseKey),
				// Enable automatic updates.
				onpremise.WithAutoUpdate(false),

				// Set the frequency in minutes that the pipeline should
				// check for updates to data files. A recommended
				// polling interval in a production environment is
				// around 30 minutes.
				onpremise.WithPollingInterval(10),
				// Set the max amount of time in seconds that should be
				// added to the polling interval. This is useful in datacenter
				// applications where multiple instances may be polling for
				// updates at the same time. A recommended amount in production
				// environments is 600 seconds.
				onpremise.WithRandomization(10),
				// Enable update on startup, the auto update system
				// will be used to check for an update before the
				// device detection engine is created.
				onpremise.WithUpdateOnStart(false),
				onpremise.WithFileWatch(true),
			)

			if err != nil {
				log.Fatalf("Failed to create engine: %v", err)
			}

			//Process evidence
			runReloadFromFileSub(pl, params.EvidenceYaml, params.DataFile)

			pl.Stop()

			return nil
		},
	)
}
