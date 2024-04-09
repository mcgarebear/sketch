package main

import (
	"fmt"
	"image/gif"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mcgarebear/sketch/cmd"

	"github.com/kelseyhightower/envconfig"
)

// Environment variable prefix for parsing application configuration
// data
const envconfigKey = "SKETCH"

// parseEnv returns the parsed structure.
func parseEnv() (*cmd.SketchConfig, error) {
	var config cmd.SketchConfig
	if err := envconfig.Process(envconfigKey, &config); err != nil {
		return nil, err
	}

	const defaultShader = ".:*o&8@#"
	if config.Shader == "" {
		config.Shader = defaultShader
	}
	config.Shader = config.Shader

	if config.Path == "" {
		return nil, fmt.Errorf("Environment variable `PATH` not provided.")
	}

	return &config, nil
}

// parseGif returns the concrete gif implementation, given a runtime configuration.
func parseGif(config *cmd.SketchConfig) (*gif.GIF, error) {
	img, err := os.Open(config.Path)
	if err != nil {
		return nil, err
	}

	gif, err := gif.DecodeAll(img)
	if err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to open image at path: " + config.Path +
			"; " + err.Error())
	}
	return gif, nil
}

func main() {
	config, err := parseEnv()
	if err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to parse environment. " + err.Error())
	}
	shaderLen := utf8.RuneCount([]byte(config.Shader))

	gif, err := parseGif(config)
	if err != nil {
		envconfig.Usagef(envconfigKey, &config, os.Stderr, envconfig.DefaultTableFormat)
		log.Fatal("Failed to parse image. " + err.Error())
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
					intensity := float32(red&0xFF)*0.2126 +
						float32(blue&0xFF)*0.7152 +
						float32(green&0xFF)*0.0722
					shaderIdx := int(intensity) % shaderLen
					fmt.Fprintf(&imageRasterized, "\x1b[38;2;%d;%d;%dm%s",
						red&0xFF, blue&0xFF, green&0xFF, string(config.Shader[shaderIdx]))
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
