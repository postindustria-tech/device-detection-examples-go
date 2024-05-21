package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/onpremise"
)

type ExampleParams struct {
	LicenseKey   string
	Product      string
	DataFile     string
	EvidenceYaml string
}

type ExampleFunc func(params ExampleParams) error

func RunExample(exampleFunc ExampleFunc) {
	licenseKey := os.Getenv("LICENSE_KEY")
	if licenseKey == "" {
		licenseKey = os.Getenv("DEVICE_DETECTION_KEY")
	}

	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "51Degrees-LiteV4.1.hash"
	}

	evidenceYaml := os.Getenv("EVIDENCE_YAML")
	if evidenceYaml == "" {
		evidenceYaml = "20000 Evidence Records.yml"
	}

	params := ExampleParams{
		LicenseKey:   licenseKey,
		DataFile:     dataFile,
		EvidenceYaml: evidenceYaml,
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

func ConverToEvidence(values map[string]string) []onpremise.Evidence {
	evidence := make([]onpremise.Evidence, len(values))
	for k, v := range values {
		strSplit := strings.SplitN(k, ".", 2)
		var prefix dd.EvidencePrefix
		switch strSplit[0] {
		case "header":
			prefix = dd.HttpHeaderString
		case "query":
			prefix = dd.HttpEvidenceQuery
		}

		evidence = append(evidence,
			onpremise.Evidence{
				Prefix: prefix,
				Key:    strSplit[1],
				Value:  v,
			})
	}
	return evidence
}

type Config struct {
	LicenseKeys  string `json:"licenseKeys"`
	DataFile     string `json:"dataFile"`
	EvidenceYaml string `json:"evidenceYaml"`
}

func LoadEnvFile(path ...string) {
	openPath := "../env.json"
	if len(path) > 0 {
		openPath = path[0]
	}

	file, err := os.Open(openPath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	// Set all defined config variables into the environment
	if _, pres := os.LookupEnv("LICENSE_KEY"); !pres {
		os.Setenv("LICENSE_KEY", config.LicenseKeys)
	}
	if _, pres := os.LookupEnv("DATA_FILE"); !pres {
		os.Setenv("DATA_FILE", config.DataFile)
	}
	if _, pres := os.LookupEnv("EVIDENCE_YAML"); !pres {
		os.Setenv("EVIDENCE_YAML", config.EvidenceYaml)
	}
}
