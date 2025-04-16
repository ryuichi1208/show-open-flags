//go:build !linux
// +build !linux

package main

import (
	"fmt"
	"os"
	"runtime"
)

// Stub function for non-Linux platforms
func runImpl() {
	// This should not be reached as the main.go already checks OS compatibility
	// but just in case
	fmt.Println("Error: This tool only supports Linux operating systems with /proc filesystem")

	// Special message for BSD systems
	bsdSystems := []string{"freebsd", "openbsd", "netbsd", "dragonfly"}
	for _, bsd := range bsdSystems {
		if runtime.GOOS == bsd {
			fmt.Println("BSD support is planned but not yet implemented")
			break
		}
	}

	fmt.Println("Current OS:", runtime.GOOS)
	os.Exit(1)
}
