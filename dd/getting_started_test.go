package dd_test

/*
This example illustrates how to perform simple device detections on given
User-Agent strings.
*/

import (
	"fmt"
	"log"

	"github.com/51Degrees/device-detection-go/dd"
)

// function match performs a match on an input User-Agent string and determine
// if the device is a mobile device. Returns output string.
func match(
	results *dd.ResultsHash,
	ua string) string {
	// Perform detection
	err := results.MatchUserAgent(ua)
	if err != nil {
		log.Fatalln("ERROR: Failed to perform detection on a User-Agent.")
	}

	propertyName := "IsMobile"

	// Get the values in string
	value, _, err := results.ValuesString(
		propertyName,
		100,
		",")
	if err != nil {
		log.Fatalln("ERROR: Failed to get results values")
	}

	hasValues, err := results.HasValues(propertyName)
	if err != nil {
		log.Fatalf(
			"ERROR: Failed to check if a matched value exists for property "+
				"%s.\n", propertyName)
	}

	if !hasValues {
		log.Printf("Property %s does not have a matched value.\n", propertyName)
	}

	return fmt.Sprintf("\tIsMobile: %s\n", value)
}

func runGettingStarted(perf dd.PerformanceProfile) string {
	// Initialise manager
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash(perf)
	filePath := getFilePath([]string{liteDataFile})

	err := dd.InitManagerFromFile(
		manager,
		*config,
		"",
		filePath)
	if err != nil {
		log.Fatalln("ERROR: Failed to initialize resource manager.")
	}

	// Make sure manager object will be freed after the function execution
	defer manager.Free()

	// Create results
	results := dd.NewResultsHash(manager, 1, 0)

	// Make sure results object is freed after function execution.
	defer results.Free()

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
	actual := fmt.Sprintf("Mobile User-Agent: %s\n", uaMobile)
	actual += match(results, uaMobile)

	// Perform detection on desktop User-Agent
	actual += fmt.Sprintf("\nDesktop User-Agent: %s\n", uaDesktop)
	actual += match(results, uaDesktop)

	// Perform detection on MediaHub User-Agent
	actual += fmt.Sprintf("\nMediaHub User-Agent: %s\n", uaMediaHub)
	actual += match(results, uaMediaHub)

	// Expected output
	expected := "Mobile User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 Safari/9537.53\n"
	expected += "\tIsMobile: True\n"
	expected += "\n"
	expected += "Desktop User-Agent: Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0\n"
	expected += "\tIsMobile: False\n"
	expected += "\n"
	expected += "MediaHub User-Agent: Mozilla/5.0 (Linux; Android 4.4.2; X7 Quad Core Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Safari/537.36\n"
	expected += "\tIsMobile: False\n"
	if actual != expected {
		log.Fatalln("Output is not as expected.")
	}
	return actual
}

func Example_getting_started() {
	// performExample(dd.Default, runGettingStarted)
	performExample(dd.BalancedTemp, runGettingStarted)
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
