package main

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"log"
)

func main() {

	RunExample(
		func(params ExampleParams) error {
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
				onpremise.WithProduct(params.Product),
				// Enable automatic updates.
				onpremise.WithAutoUpdate(true),
				// Passing a custom update URL
				//onpremise.WithDataUpdateUrl(
				//	"https://myprovider.com/1.tar.gz",
				//),
				// Enable update on startup, the auto update system
				// will be used to check for an update before the
				// device detection engine is created.
				onpremise.WithUpdateOnStart(true),

				// Set the frequency in minutes that the engine should
				// check for updates to data files. A recommended
				// polling interval in a production environment is
				// around 30 minutes.
				onpremise.WithPollingInterval(1),

				//set custom polling interval
				// randomize polling interval, it will append random number of seconds to polling interval
				// in case you have multiple instances of the engine running and want to avoid polling at the exact same time
				onpremise.WithRandomization(60),

				//File watching is enabled by default any external changes to the data file will be picked up and reloaded
				//onpremise.WithFileWatch(true),

				//TempDataCopy is enabled by default, it will copy the data file to a temporary location before loading
				//onpremise.WithTempDataCopy(true),
				//SetTempDataDir is used to set the directory where the temporary data file will be copied
				//default is os.TempDir() aka temporary directory of the OS
				//onpremise.SetTempDataDir("/tmp"),
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
