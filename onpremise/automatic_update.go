package main

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"log"
)

func main() {
	//Create config
	config := dd.NewConfigHash(dd.Balanced)

	//Create on-premise engine
	pl, err := onpremise.New(
		config,
		//Provide your own file URL
		//it can be compressed as gz or raw, engine will handle it
		onpremise.WithDataUpdateUrl(
			"https://myprovider.com/1.tar.gz",
		),
		//set custom polling interval
		onpremise.WithPollingInterval(2),
		// randomize polling interval, it will append random number of seconds to polling interval
		// in case you have multiple instances of the engine running and want to avoid polling at the exact same time
		onpremise.WithRandomization(60),

		// In case you want to use 51degrees file provider you just set licence key and product
		onpremise.WithProduct("Hash"),
		onpremise.WithLicenceKey("YOUR_LICENCE_KEY"),
		//	in case you want to disable auto update
		onpremise.WithAutoUpdate(false),
		// in case you want to update file from URL on start of the engine otherwise it will be pulled
		// on next polling interval
		onpremise.WithUpdateOnStart(true),
	)

	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	//Process evidence
	resultsHash, err := pl.Process(
		[]onpremise.Evidence{
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua-Arch",
				Value:  "x86",
			},
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua-Model",
				Value:  "Intel",
			},
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua-Mobile",
				Value:  "?0",
			},
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua-Platform",
				Value:  "Windows",
			},
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua-Platform-Version",
				Value:  "10.0",
			},
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua-Full-Version-List",
				Value:  "58.0.3029.110",
			},
			{
				Prefix: dd.HttpHeaderString,
				Key:    "Sec-Ch-Ua",
				Value:  `"\"Chromium\";v=\"91.0.4472.124\";a=\"x86\";p=\"Windows\";rv=\"91.0\""`,
			},
		},
	)
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
	<-make(chan struct{})

	pl.Stop()
}
