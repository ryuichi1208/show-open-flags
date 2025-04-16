//go:build linux
// +build linux

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type FileInfo struct {
	fileName string
	fd       int64
	mode     fs.FileMode
}

func getFDList(fileName string) ([]*FileInfo, error) {
	files, err := os.ReadDir(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var fileList []*FileInfo
	for _, file := range files {
		info, err := os.Lstat(fmt.Sprintf("%s/%s", fileName, file.Name()))
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			realPath, _ := os.Readlink(fmt.Sprintf("%s/%s", fileName, file.Name()))
			if realPath[0:1] == "/" {
				i, err := strconv.Atoi(file.Name())
				if err != nil {
					return nil, err
				}
				fi := &FileInfo{
					fileName: realPath,
					fd:       int64(i),
					mode:     info.Mode(),
				}
				fileList = append(fileList, fi)
			}

		} else {
			continue
		}
	}

	return fileList, nil
}

func readFDInfo(fileName string) []byte {
	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 64)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if strings.Contains(string(line), "flags") {
			return line
		}
	}

	return []byte("")
}

const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
	// other
	O_ASYNC     int = syscall.O_ASYNC
	O_CLOEXEC   int = syscall.O_CLOEXEC
	O_DIRECT    int = 00040000 // Defined as a constant as it may only be defined on Linux
	O_DIRECTORY int = syscall.O_DIRECTORY
	O_DSYNC     int = syscall.O_DSYNC
	O_FSYNC     int = syscall.O_FSYNC
	O_LARGEFILE int = 00100000 // Defined as a constant for the same reason
	O_NDELAY    int = syscall.O_NDELAY
	O_RSYNC     int = 00040000 // Defined as a constant for the same reason
)

func checkFlags(hex int64) []string {
	var fs []string
	if hex&int64(O_RDONLY) > 0 {
		fs = append(fs, "O_RDONLY")
	}
	if hex&int64(O_WRONLY) > 0 {
		fs = append(fs, "O_WRONLY")
	}
	if hex&int64(O_RDWR) > 0 {
		fs = append(fs, "O_RDWR")
	}
	if hex&int64(O_APPEND) > 0 {
		fs = append(fs, "O_APPEND")
	}
	if hex&int64(O_CREATE) > 0 {
		fs = append(fs, "O_CREATE")
	}
	if hex&int64(O_EXCL) > 0 {
		fs = append(fs, "O_EXCL")
	}
	if hex&int64(O_SYNC) > 0 {
		fs = append(fs, "O_SYNC")
	}
	if hex&int64(O_TRUNC) > 0 {
		fs = append(fs, "O_TRUNC")
	}
	if hex&int64(O_CREATE) > 0 {
		fs = append(fs, "O_CREATE")
	}
	// if hex&int64(O_ACCMODE) > 0 {
	// 	fs = append(fs, "O_ACCMODE")
	// }
	if hex&int64(O_ASYNC) > 0 {
		fs = append(fs, "O_ASYNC")
	}
	if hex&int64(O_CLOEXEC) > 0 {
		fs = append(fs, "O_CLOEXEC")
	}
	if hex&int64(O_DIRECT) > 0 {
		fs = append(fs, "O_DIRECT")
	}
	if hex&int64(O_DIRECTORY) > 0 {
		fs = append(fs, "O_DIRECTORY")
	}
	if hex&int64(O_DSYNC) > 0 {
		fs = append(fs, "O_DSYNC")
	}
	if hex&int64(O_FSYNC) > 0 {
		fs = append(fs, "O_FSYNC")
	}
	if hex&int64(O_LARGEFILE) > 0 {
		fs = append(fs, "O_LARGEFILE")
	}
	if hex&int64(O_NDELAY) > 0 {
		fs = append(fs, "O_NDELAY")
	}
	if hex&int64(O_RSYNC) > 0 {
		fs = append(fs, "O_RSYNC")
	}

	return fs
}

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
