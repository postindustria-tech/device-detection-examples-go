package dd_test

/*
This example illustrates how to perform simple device detections on given
User-Agent strings.
*/

import (
	"fmt"

	"github.com/51Degrees/device-detection-go/dd"
)

// function match performs a match on an input User-Agent string and determine
// if the device is a mobile device.
func match(
	results dd.ResultsHash,
	ua string) {
	// Perform detection
	err := results.MatchUserAgent(ua)
	if err != nil {
		panic(err)
	}

	// Get the values in string
	value, _, err := results.ValuesString(
		"IsMobile",
		100,
		",")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\tIsMobile: %s\n", value)
}

func Example_getting_started() {
	// Initialise manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash()
	filePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	err := dd.InitManagerFromFile(
		manager,
		config,
		"",
		filePath)
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

	// Create results
	results, err := dd.NewResultsHash(
		manager,
		1,
		0)
	if err != nil {
		panic(err)
	}

	// Make sure results object is freed after function execution.
	defer func() {
		err = results.Free()
		if err != nil {
			panic(err)
		}
	}()

	// User-Agent string of an iPhone mobile device.
	const uaMobile = "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) " +
		"AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 " +
		"Safari/9537.53"

	// User-Agent string of Firefox Web browser version 41 on desktop.
	const uaDesktop = "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) " +
		"Gecko/20100101 Firefox/41.0"

	// User-Agent string of a MediaHub device.
	const uaMediaHub = "Mozilla/5.0 (Linux; Android 4.4.2; X7 Quad Core " +
		"Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 " +
		"Chrome/30.0.0.0 Safari/537.36"

	// Perform detection on mobile User-Agent
	fmt.Printf("Mobile User-Agent: %s\n", uaMobile)
	match(results, uaMobile)

	// Perform detection on desktop User-Agent
	fmt.Printf("\nDesktop User-Agent: %s\n", uaDesktop)
	match(results, uaDesktop)

	// Perform detection on MediaHub User-Agent
	fmt.Printf("\nMediaHub User-Agent: %s\n", uaMediaHub)
	match(results, uaMediaHub)

	// Output:
	// Mobile User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 Safari/9537.53
	// 	IsMobile: True
	//
	// Desktop User-Agent: Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0
	// 	IsMobile: False
	//
	// MediaHub User-Agent: Mozilla/5.0 (Linux; Android 4.4.2; X7 Quad Core Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Safari/537.36
	// 	IsMobile: False
}
