package main

import (
	"image/gif"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type sketchConfig struct {
	Color  bool   `envconfig:"color"`
	Path   string `enconfig:"path"`
	Frame  int    `enconfig:"frame"`
	Shader string `envconfig:"shader"`
}

func main() {

	const envconfigKey = "SKETCH"
	var config sketchConfig
	if err := envconfig.Process(envconfigKey, &config); err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to process environment variables: " + err.Error())
	}

	if config.Path == "" {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to open image at path: " + config.Path +
			"; Path not provided.")
	}

	img, err := os.Open(config.Path)
	if err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to open image at path: " + config.Path +
			"; " + err.Error())
	}

	gif, err := gif.DecodeAll(img)
	if err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to open image at path: " + config.Path +
			"; " + err.Error())
	}

	log.Println(gif)

}
