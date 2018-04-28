package lwc

import (
	"fmt"
	"os"
	"strings"
)

func Run(version string, date string) {
	// Read command-line args
	config := NewConfig(os.Args)

	switch {
	case config.Version:
		// Print version and exit
		fmt.Printf("lwc %s\n", version)
		if date != "" {
			fmt.Printf("Built %s\n", strings.Replace(date, "_", " ", -1))
		}
	case config.Help:
		// Print usage and exit
		config.PrintUsage()
	default:
		// Process input
		ProcessFiles(config)
	}
}
