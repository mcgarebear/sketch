package main

import (
	"fmt"
	"image/gif"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type sketchConfig struct {
	Color  bool   `envconfig:"color"`
	Path   string `enconfig:"path"`
	Frame  int    `enconfig:"frame"`
	Shader string `envconfig:"shader"`
}

func main() {

	// parse configuration from environment
	const envconfigKey = "SKETCH"
	var config sketchConfig
	if err := envconfig.Process(envconfigKey, &config); err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to process environment variables: " + err.Error())
	}

	// verify path and attempt to open file
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

	// decode open file into gif
	gif, err := gif.DecodeAll(img)
	if err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to open image at path: " + config.Path +
			"; " + err.Error())
	}
	log.Println("gif", gif)
	log.Println("gif.LoopCount", gif.LoopCount)
	log.Println("gif.Config", gif.Config)
	log.Println("gif.BackgroundIndex", gif.BackgroundIndex)

	// for each image in the gif
	numImages := len(gif.Image)
	log.Println("numImages", numImages)
	for idx := 0; idx < numImages; idx++ {
		log.Println("idx", idx)
		log.Println("gif.Image", gif.Image[idx])
		log.Println("gif.Delay", gif.Delay[idx])
		log.Println("gif.Disposal", gif.Disposal[idx])

		// iterate through each pixel in the image, starting from the top left
		// and moving to the bottom right. For each pixel, index into the color
		// pallette to determine the pixel's color value.
		image := gif.Image[idx]
		for row := 0; row < gif.Config.Height; row++ {
			for col := 0; col < gif.Config.Width; col++ {
				pixelIdx := (row-image.Rect.Min.Y)*image.Stride + (col - image.Rect.Min.X)
				color := image.Palette[image.Pix[pixelIdx]]
				log.Println("x", col, "y", row, "color", color)
				r, g, b, _ := color.RGBA()
				fmt.Printf("\x1b[38;5;%dmx", r+g+b%8)
			}
			fmt.Printf("\n")
		}
		// clear color, reset cursor
		fmt.Printf("\x1b[0m")
		time.Sleep(time.Second / 100 * time.Duration(gif.Delay[idx]))
		if idx < numImages {
			fmt.Printf("\x1b[%dF", gif.Config.Height)
		}
	}

}
