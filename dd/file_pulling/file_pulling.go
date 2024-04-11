package main

import (
	"github.com/51Degrees/device-detection-go/v4/dd"
	"github.com/51Degrees/device-detection-go/v4/pipeline"
	"log"
)

func main() {
	manager := dd.NewResourceManager()
	config := dd.NewConfigHash(dd.Balanced)

	pl, err := pipeline.New(
		manager,
		config,
		pipeline.WithDataUpdateUrl(
			"https://distributor.51degrees.com/api/v2/download?LicenseKeys=FLVAAAS7TAAASA6CEJAEAW22HDBQMX6CTF6KJ4DHQ2GGFDE4MQDMC99R5K65Y98HPTQC7MS38V7WN688MBL28GB&Type=HashV41&Download=True&Product=V4TAC",
			10000,
		),
	)

	if err != nil {
		log.Fatalf("Failed to create pipeline: %v", err)
	}

	err = pl.Run()
	if err != nil {
		log.Fatalf("Failed to run pipeline: %v", err)
	}
	<-make(chan struct{})

	defer manager.Free()

}
