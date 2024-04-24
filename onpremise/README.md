## onpremise Engine

This example demonstrates how to use the on-premise engine to detect devices.



### Running the example

#### Create config
```go
    config := dd.NewConfigHash(dd.Balanced)
```

#### Create engine
```go
    e, err := New(
                config,
                WithDataUpdateUrl("datafileUrl.com/myFile.gz", 2000),
				WithDataFile("51Degrees-LiteV4.1.hash"),
         )
```

#### Process
```go
resultsHash, err := e.Process(
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







