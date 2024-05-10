package main

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"os"
)

type ExampleParams struct {
	LicenceKey string
	Product    string
	DataFile   string
}

type ExampleFunc func(params ExampleParams) error

func RunExample(exampleFunc ExampleFunc) {
	params := ExampleParams{
		LicenceKey: os.Getenv("LICENCE_KEY"),
		DataFile:   "51Degrees-LiteV4.1.hash",
	}

	err := exampleFunc(params)
	if err != nil {
		panic(err)
	}
}

var ExampleEvidence = []onpremise.Evidence{
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
}
