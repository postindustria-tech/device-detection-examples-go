# 51Degrees Device Detection Engines

![51Degrees](https://51degrees.com/DesktopModules/FiftyOne/Distributor/Logo.ashx?utm_source=github&utm_medium=repository&utm_content=readme_main&utm_campaign=go-open-source "Data rewards the curious") **Examples for Device Detection in Go**

## Introduction

This repository contains examples of how to use module [device-detection-go](https://github.com/51degrees/device-detection-go)

## Pre-requisites
To run these examples you will need a data file and example evidence for some of the tests.  To fetch these assets please run:

```
pwsh ci/fetch-assets.ps1 .
```

or alternatively you can download them from [device-detection-data](https://github.com/51Degrees/device-detection-data) repo (the links are below) and put in the root of this repository. 

- [51Degrees-LiteV4.1.hash](https://github.com/51Degrees/device-detection-data/blob/main/51Degrees-LiteV4.1.hash)
- [20000 User Agents.csv](https://github.com/51Degrees/device-detection-data/blob/main/20000%20User%20Agents.csv)

## Examples

**NOTE**: `device-detection-examples-go` references `device-detection-go` as a dependency in `go.mod`.  No additional actions should be required - the module will be downloaded and built when you do `go run`, `go test`, or `go build` explicitly for any example.  

The examples are grouped into Go testable examples and Go app.

- All examples under folder `dd` are testable examples and suffixed with `*_test.go` at the end. Expected output of these examples are in the comment section of each example (i.e.  // Output:...). This example can be run by `go test`. For example, to run the example `getting_started_test.go` we can use `go test`, passing the `Example_*` method name of the example which acts as the main entry point. The following command is run in the `dd` folder.
```
go test -run Example_getting_started
```
- Example under the `web` folder is a Go web application that can be run using `go run`.
```
go run web_integration.go
```

Below is a table that describes the examples:

|Example|Description|
|-------|-----------|
|dd/getting_started_test.go|A simple example that shows how to initialize a resource manager and perform device detection on User-Agent strings.|
|dd/match_device_id_test.go|A simple example that shows how to perform device detection using Device Id.|
|dd/match_metrics_test.go|A simple example that shows how to access match metrics.|
|dd/offline_processing_test.go|An example that shows how to process through User-Agents stored in a file, and output detection results and metrics to a local file for further evaluation. Output file is `./device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.processed.csv`|
|dd/performance_test.go|An example perform performance benchmarking of our device detection solution and output the benchmark to a report file. Output file is `performance_report.log` in the working directory.|
|dd/reload_from_file_test.go|An example that demonstrates how a data file can be reloaded while serving device detection requests.|
|dd/reload_from_memory_test.go|To be implemented|
|dd/strongly_typed_test.go|To be implemented|
|web/web_integration.go|An example of how `device-detection-go` can be used in a web application.|
|uach/uach.go|An example of how `User Agent Client Hints (UACH)` can be requested by the `Device Detection` engine and how they can be used as evidence to perform a detection. Please also read the comment at the top of the example file `uach.go` which also provides a greater details on usage of UACH with `Device Detection` engine.|

## Run examples

- Navigate to `dd` folder. All examples here are testable and can be run as:
```
go test -run [Name of Example_* method of the target example]
```
- Navigate to `web` folder. This is a web app and it can be run as:
```
go run web_integration.go
```
- Navigate to `uach` folder. This is a web app and it can be run as:
```
go run uach.go
```

For futher details of how to run each example, please read more at the comment section located at the top of each example file.