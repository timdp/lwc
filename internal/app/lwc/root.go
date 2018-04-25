package lwc

import (
	"fmt"
	"os"
)

func Run(version string) {
	// Read command-line args
	config := BuildConfig(os.Args)

	switch {
	case config.Version:
		// Print version and exit
		fmt.Printf("lwc %s\n", version)
	case config.Help:
		// Print usage and exit
		config.PrintUsage()
	default:
		// Process input
		processors := BuildProcessors(&config)
		ProcessFiles(&config, processors)
	}
}
