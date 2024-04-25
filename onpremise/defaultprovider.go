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
	engine, err := onpremise.New(
		config,
		onpremise.SetLicenceKey("123"),
		onpremise.SetProduct("MyProduct"),
	)

	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	//Process evidence
	resultsHash, err := engine.Process(
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

	engine.Stop()
}
