package main

import (
	"fmt"
	"os"
	"runtime"
)

// IsSupportedOS checks if the current OS is supported
func IsSupportedOS() (bool, string) {
	// Linux is fully supported
	if runtime.GOOS == "linux" {
		return true, ""
	}

	// Check for BSD variants
	bsdSystems := []string{"freebsd", "openbsd", "netbsd", "dragonfly"}
	for _, bsd := range bsdSystems {
		if runtime.GOOS == bsd {
			return false, "BSD support is planned but not yet implemented"
		}
	}

	// Other systems are not supported
	return false, fmt.Sprintf("This tool does not support %s operating systems", runtime.GOOS)
}

func main() {
	// Check if the OS is supported
	supported, message := IsSupportedOS()

	if !supported {
		fmt.Println("Error: This tool only supports Linux operating systems with /proc filesystem")
		if message != "" {
			fmt.Println(message)
		}
		fmt.Println("Current OS:", runtime.GOOS)
		os.Exit(1)
	}

	// If arguments are insufficient, show usage
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<PID>")
		fmt.Println("Please provide a process ID (PID)")
		os.Exit(1)
	}

	// Call the OS-specific implementation
	runImpl()
}

