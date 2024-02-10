package main

import (
	"fmt"
	"image/gif"
	"log"
	"os"
	"strings"
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

	// hide cursor, best effort to restore prompt
	fmt.Printf("\x1b?25l")
	defer fmt.Printf("\x1b?25h")
	defer fmt.Printf("\x1b[0m")

	// for each image in the gif
	numImages := len(gif.Image)
	prevHeight := 0
	for idx := 0; idx < numImages; idx++ {
		var imageRasterized strings.Builder
		// iterate through each pixel in the image, starting from the top left
		// and moving to the bottom right. For each pixel, index into the color
		// pallette to determine the pixel's color value.
		image := gif.Image[idx]
		height := image.Bounds().Dy()
		width := image.Bounds().Dx()
		for row := 0; row < height; row++ {
			for col := 0; col < width; col++ {
				red, green, blue, alpha := image.At(col, row).RGBA()
				if alpha > 0 {
					fmt.Fprintf(&imageRasterized, "\x1b[38;2;%d;%d;%dmx",
						red&0xFF, blue&0xFF, green&0xFF)
				} else {
					fmt.Fprintf(&imageRasterized, "\x1b[0m ")
				}

			}
			// handle gifs with varying sized frames:
			//   clear from cursor to EOL. newline for next row.
			fmt.Fprintf(&imageRasterized, "\x1b[0J")
			fmt.Fprintf(&imageRasterized, "\n")
		}
		// handle gifs with varying sized frames:
		//   clear any additional rows from the previous frame by clearing the line
		for deltaHeight := prevHeight - height; deltaHeight > 0; deltaHeight-- {
			fmt.Fprintf(&imageRasterized, "\x1b[2K\n")
		}
		// move cursor back to original position (except for last frame)
		if (idx + 1) != numImages {
			fmt.Fprintf(&imageRasterized, "\x1b[1;1H")
		}
		prevHeight = height
		// draw, then sleep for animation
		fmt.Printf(imageRasterized.String())
		time.Sleep((time.Second / 100) * time.Duration(gif.Delay[idx]))
	}

}
