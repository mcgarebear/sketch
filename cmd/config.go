package cmd

// SketchConfig represents the configurable attributes of the `sketch`
// executable.
type SketchConfig struct {
	Path   string `enconfig:"path"`
	Shader string `envconfig:"shader"`
}
