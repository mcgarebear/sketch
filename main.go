package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type sketchConfig struct {
	Color bool `envconfig:"color"`
	Path string `enconfig:"path"`
	Frame int `enconfig:"frame"`
}

func main() {

	const envconfigKey = "SKETCH"
	var config sketchConfig
	if err := envconfig.Process(envconfigKey, &config); err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to process environment variables: " + err.Error())
	}

	log.Println(config)
}
