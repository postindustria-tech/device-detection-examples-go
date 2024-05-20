package common

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
	"os"
)

type ExampleParams struct {
	LicenseKey string
	Product    string
	DataFile   string
}

type ExampleFunc func(params ExampleParams) error

func RunExample(exampleFunc ExampleFunc) {
	licenseKey := os.Getenv("LICENSE_KEY")
	if licenseKey == "" {
		licenseKey = os.Getenv("DEVICE_DETECTION_KEY")
	}

	params := ExampleParams{
		LicenseKey: licenseKey,
		DataFile:   "51Degrees-LiteV4.1.hash",
	}

	err := exampleFunc(params)
	if err != nil {
		panic(err)
	}
}

var ExampleEvidence1 = []onpremise.Evidence{
	{Prefix: dd.HttpHeaderString, Key: "Sec-Ch-Ua", Value: "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\""},
	{Prefix: dd.HttpHeaderString, Key: "Sec-Ch-Ua-Full-Version-List", Value: "\"Chromium\";v=\"124.0.6367.208\", \"Google Chrome\";v=\"124.0.6367.208\", \"Not-A.Brand\";v=\"99.0.0.0\""},
	{Prefix: dd.HttpHeaderString, Key: "Sec-Ch-Ua-Mobile", Value: "?0"},
	{Prefix: dd.HttpHeaderString, Key: "Sec-Ch-Ua-Platform", Value: "\"macOS\""},
	{Prefix: dd.HttpHeaderString, Key: "Sec-Ch-Ua-Platform-Version", Value: "\"14.4.1\""},
	{Prefix: dd.HttpHeaderString, Key: "User-Agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"},
}

var ExampleEvidence2 = []onpremise.Evidence{
	{Prefix: dd.HttpHeaderString, Key: "User-Agent", Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"},
}
