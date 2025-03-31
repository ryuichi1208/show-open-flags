//go:build !linux
// +build !linux

package main

import (
	"fmt"
	"os"
)

// Stub function for non-Linux platforms
func runLinux() {
	fmt.Println("Error: This tool only supports Linux operating systems with /proc filesystem")
	os.Exit(1)
}
