package main

/*
This example illustrates how to perform simple device detections on given
User-Agent strings.
*/

import (
	"html/template"
	"log"
	"net/http"

	"github.com/51Degrees/device-detection-go/dd"
)

// Properties required for a response page.
type Page struct {
	BrowserName       string
	ScreenPixelsWidth string
}

var manager *dd.ResourceManager
var config *dd.ConfigHash

// Template for the response HTML page.
var templ1 = `<!DOCTYPE HTML>
<html>
  <head>
    <meta charset="utf-8">
    <title>Web Integration Example</title>
  </head>
  <body>
    <p id=browsername>Browser: <b>{{.BrowserName}}</b></p>
    <p id=screenpixelswidth>Screen Pixels Width: <b>{{.ScreenPixelsWidth}}</b></p>
  </body>
</html>`

// function match performs a match on an input User-Agent string and determine
// if the device is a mobile device.
func match(
	results *dd.ResultsHash,
	ua string) {
	// Perform detection
	err := results.MatchUserAgent(ua)
	if err != nil {
		log.Fatalf("ERROR: Failed to perform detection on User-Agent \"%s\".\n", ua)
	}
}

// function getValue return a value results for a property
func getValue(
	results *dd.ResultsHash,
	propertyName string) string {
	// Get the values in string
	value, _, err := results.ValuesString(
		propertyName,
		100,
		",")
	if err != nil {
		log.Fatalln("ERROR: Failed to get results values string.")
	}

	hasValues, err := results.HasValues(propertyName)
	if err != nil {
		log.Fatalf(
			"ERROR: Failed to check if a matched value exists for property "+
				"%s.\n", propertyName)
	}

	if !hasValues {
		log.Printf("Property %s does not have a matched value.\n", propertyName)
		return ""
	}

	return value
}

// Handler for web request
func handler(w http.ResponseWriter, r *http.Request) {
	// Create results
	results, err := dd.NewResultsHash(
		manager,
		1,
		0)
	if err != nil {
		log.Fatalln("ERROR: Failed to create new results.")
	}

	// Make sure results object is freed after function execution.
	defer func() {
		err = results.Free()
		if err != nil {
			log.Fatalln("ERROR: Failed to free results.")
		}
	}()

	// Perform detection on mobile User-Agent
	match(results, r.UserAgent())
	browserName := getValue(results, "BrowserName")
	screenPixelWidth := getValue(results, "ScreenPixelsWidth")
	p := &Page{
		browserName,
		screenPixelWidth,
	}

	// Construct the template
	t := template.Must(template.New("dd").Parse(templ1))
	// Return the constructed template in a response
	t.Execute(w, p)
}

func main() {
	// Initialise manager
	manager = dd.NewResourceManager()
	config = dd.NewConfigHash(dd.Balanced)
	filePath := "../device-detection-go/dd/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	err := dd.InitManagerFromFile(
		manager,
		*config,
		"",
		filePath)
	if err != nil {
		log.Fatalln("ERROR: Failed to initialize resource manager.")
	}

	// Make sure manager object will be freed after the function execution
	defer func() {
		err := manager.Free()
		if err != nil {
			log.Fatalln("ERROR: Failed to free resource manager.")
		}
	}()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
