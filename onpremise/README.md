## onpremise Engine

This example demonstrates how to use the 51degrees onpremise engine to detect devices.



### Running the example

#### Create config
```go
    config := dd.NewConfigHash(dd.Balanced)
```

#### Create engine
```go
    engine, err := New(
                config,
                WithDataUpdateUrl("datafileUrl.com/myFile.gz", 2000),
				WithDataFile("51Degrees-LiteV4.1.hash"),
         )
```

#### Process
```go
resultsHash, err := engine.Process(
        []Evidence{
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
		}
)

```

#### Get values
```go

browser, err := resultsHash.ValuesString("BrowserName", ",")
	if err != nil {
		log.Fatalf("Failed to get BrowserName: %v", err)
	}
```

#### Dont forget to free memory
```go
 defer resultsHash.Free()
```

### Options

#### WithDataUpdateUrl Provides existing datafile
* path - path to the datafile in case you stored it locally
```go
    WithDataFile(path string) EngineOptions
```

#### WithDataUpdateUrl Provides datafile update url
use this in case you have your own datafile source url and want to update it.
you can see example in updatedatafile.go
* url - url to the datafile
* interval - interval in seconds for fetching the datafile
```go
    WithDataUpdateUrl(url string, interval int) EngineOptions
```

####SetLicenceKey
in case you use default 51degrees datafile provider you need to set licence key and product
you can see example in defaultprovider.go
```go
   SetLicenceKey(key string) EngineOptions
   SetProduct(product string) EngineOptions
```

#### ToggleLogger Enables or disables logger
* enable - true or false
```go
    ToggleLogger(enabled bool) EngineOptions
```

####WithCustomLogger Provides custom logger
* logger - custom logger
  * Logger muster implement LogWriter interface
```go
    WithCustomLogger(logger LogWriter) EngineOptions
```

####SetMaxRetries
 enables you to set maximum retries for fetching datafile, default is 0, meaning infinite retries
```go
SetMaxRetries(retries int) EngineOptions
```






