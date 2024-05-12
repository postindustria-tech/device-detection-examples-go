package main

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"log"
)

func main() {

	RunExample(
		func(params ExampleParams) error {
			//... Example code
			//Create config
			config := dd.NewConfigHash(dd.Balanced)

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
				onpremise.WithLicenceKey(params.LicenceKey),
				// Enable automatic updates.
				onpremise.WithAutoUpdate(true),

				// Set the frequency in minutes that the pipeline should
				// check for updates to data files. A recommended
				// polling interval in a production environment is
				// around 30 minutes.
				onpremise.WithPollingInterval(1),
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
			)

			if err != nil {
				log.Fatalf("Failed to create engine: %v", err)
			}

			//Process evidence
			resultsHash, err := pl.Process(ExampleEvidence)
			if err != nil {
				log.Fatalf("Failed to process: %v", err)
			}
			defer resultsHash.Free()

			//Get values from results
			browser, err := resultsHash.ValuesString("BrowserName", ",")
			if err != nil {
				log.Fatalf("Failed to get BrowserName: %v", err)
			}

			log.Printf("BrowserName: %s", browser)

			deviceType, err := resultsHash.ValuesString("DeviceType", ",")
			if err != nil {
				log.Fatalf("Failed to get DeviceType: %v", err)
			}

			log.Printf("DeviceType: %s", deviceType)

			//use results and do detection
			resultsHash.Free()

			pl.Stop()

			return nil
		},
	)
}
