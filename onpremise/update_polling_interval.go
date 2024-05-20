package main

import (
	"github.com/51Degrees/device-detection-examples-go/v4/onpremise/common"
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"log"
	"time"
)

func processExampleEvidence(engine *onpremise.Engine) {
	//Process evidence
	resultsHash, _ := engine.Process(common.ExampleEvidence)
	defer resultsHash.Free()
	//Get values from results
	browser, _ := resultsHash.ValuesString("BrowserName", ",")
	deviceType, _ := resultsHash.ValuesString("DeviceType", ",")
	log.Printf("BrowserName: %s", browser)
	log.Printf("DeviceType: %s", deviceType)
}

func main() {

	common.RunExample(
		func(params common.ExampleParams) error {
			//... Example code
			//Create config
			config := dd.NewConfigHash(dd.Balanced)

			//Create on-premise engine
			engine, err := onpremise.New(
				config,

				// Path to your data file
				onpremise.WithDataFile(params.DataFile),

				// For automatic updates to work you will need to provide a license key.
				// A license key can be obtained with a subscription from https://51degrees.com/pricing
				onpremise.WithLicenseKey(params.LicenseKey),

				// Enable automatic updates.
				onpremise.WithAutoUpdate(true),

				// Set the frequency in minutes that the pipeline should
				// check for updates to data files. A recommended
				// polling interval in a production environment is
				// around 30 minutes.
				onpremise.WithPollingInterval(5),

				// Set the max amount of time in seconds that should be
				// added to the polling interval. This is useful in datacenter
				// applications where multiple instances may be polling for
				// updates at the same time. A recommended amount in production
				// environments is 600 seconds.
				onpremise.WithRandomization(1),

				// Enable update on startup, the auto update system
				// will be used to check for an update before the
				// device detection engine is created.
				onpremise.WithUpdateOnStart(false),

				// Optionally provide your own file URL
				// onpremise.WithDataUpdateUrl(""),

				// Whether a temp copy should be created
				onpremise.WithTempDataCopy(false),
			)

			if err != nil {
				log.Fatalf("Failed to create engine: %v", err)
			}

			defer engine.Stop()

			//process before file has been updated
			processExampleEvidence(engine)

			<-time.After(20 * time.Second)

			//process again after the file has presumably been updated
			processExampleEvidence(engine)
			processExampleEvidence(engine)
			processExampleEvidence(engine)

			return nil
		},
	)
}
