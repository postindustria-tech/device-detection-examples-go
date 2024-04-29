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

### Software

In order to use device-detection-examples-go the following are required:
- A C compiler that support C11 or above (Gcc on Linux, Clang on MacOS and MinGW-x64 on Windows)
- libatomic - which usually come with default Gcc, Clang installation

### Windows

If you are on Windows, make sure that:
- The path to the `MinGW-x64` `bin` folder is included in the `PATH`. By default, the path should be `C:\msys64\ucrt64\bin`
- Go environment variable `CGO_ENABLED` is set to `1` 
```
go env -w CGO_ENABLED=1
```

## Examples

**NOTE**: `device-detection-examples-go` references `device-detection-go` as a dependency in `go.mod`.  No additional actions should be required - the module will be downloaded and built when you do `go run`, `go test`, or `go build` explicitly for any example.  

- All examples under `dd` directory are console program examples and are run using `go run`.
- Example under the `web` and `uach` directories are Go web applications that can also be run using `go run`.

Below is a table that describes the examples:

|Example|Description|
|-------|-----------|
|dd/getting_started/getting_sarted.go|A simple example that shows how to initialize a resource manager and perform device detection on User-Agent strings.|
|dd/match_device_id/match_device_id.go|A simple example that shows how to perform device detection using Device Id.|
|dd/match_metrics/match_metrics.go|A simple example that shows how to access match metrics.|
|dd/offline_processing/offline_processing.go|An example that shows how to process through User-Agents stored in a file, and output detection results and metrics to a local file for further evaluation. Output file is `./device-detection-go/dd/device-detection-cxx/device-detection-data/20000 User Agents.processed.csv`|
|dd/performance/performance.go|An example perform performance benchmarking of our device detection solution and output the benchmark to a report file. Output file is `performance_report.log` in the working directory.|
|dd/reload_from_file/reload_from_file.go|An example that demonstrates how a data file can be reloaded while serving device detection requests.|
|dd/reload_from_memory/reload_from_memory.go|To be implemented|
|dd/strongly_typed/strongly_typed.go|To be implemented|
|web/web_integration.go|An example of how `device-detection-go` can be used in a web application.|
|uach/uach.go|An example of how `User Agent Client Hints (UACH)` can be requested by the `Device Detection` engine and how they can be used as evidence to perform a detection. Please also read the comment at the top of the example file `uach.go` which also provides a greater details on usage of UACH with `Device Detection` engine.|

## Run examples

- Navigate to `dd` folder. All examples here are testable and can be run as:
```
go run [example_dir/example_name].go
```
- Navigate to `web` folder. This is a web app and it can be run as:
```
go run web_integration.go
```
- Navigate to `uach` folder. This is a web app and it can be run as:
```
go run uach.go
```

For futher details of how to run each example, please read more in the comment section located at the top of each example file.