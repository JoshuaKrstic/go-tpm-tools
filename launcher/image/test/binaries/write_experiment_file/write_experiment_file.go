package main

import (
	"os"
	"path"

	"github.com/google/go-tpm-tools/launcher/launcherfile"
)

const (
	experimentDataFile = "experiment_data"
)

func main() {
	if err := os.MkdirAll(launcherfile.HostTmpPath, 0744); err != nil {
		panic(err)
	}
	experimentsFile := path.Join(launcherfile.HostTmpPath, experimentDataFile)

	data := "{\"EnableTestFeatureForImage\":true,\"NonExistantExperimentFlag\":true}"

	err := os.WriteFile(experimentsFile, []byte(data), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
