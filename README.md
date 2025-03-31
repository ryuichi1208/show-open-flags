# show-open-flags

A command-line tool to display the open(2) flags used by file descriptors in a running process. This utility helps you analyze how a process has opened its files, which can be useful for debugging, security analysis, and understanding process behavior.

## Overview

`show-open-flags` examines a process's file descriptors by reading the `/proc` filesystem and displays the open flags (such as `O_RDONLY`, `O_WRONLY`, `O_RDWR`, `O_APPEND`, etc.) for each open file. This information is not easily available through standard tools like `lsof` but can be critical for understanding process behavior.

## Platform Support

This tool is **Linux-only** and will not work on Windows, macOS, or other non-Linux platforms. The code includes:

- Build constraints that prevent compilation on non-Linux platforms
- Runtime checks that display an error message if somehow executed on a non-Linux system

This restriction exists because the tool relies heavily on the Linux `/proc` filesystem which is not available on other operating systems.

## System Requirements

- Linux-based operating system with `/proc` filesystem
- Go 1.23 or later (for building from source)

## Installation

### From Source

1. Clone the repository:

```bash
git clone https://github.com/ryuichi1208/show-open-flags.git
cd show-open-flags
```

2. Build the binary:

```bash
go build
```

3. Install (optional):

```bash
sudo cp show-open-flags /usr/local/bin/
```

## Testing

Run the included tests with:

```bash
go test -v
```

The test suite includes tests for flag parsing, file descriptor information reading, and other core functionality. Some tests are skipped due to environment dependencies or because they would require specific system paths.

## Usage

```bash
show-open-flags PID
```

Where:
- `PID` is the process ID of the process you want to analyze

## Examples

### Examining the init process (PID 1)

```bash
$ show-open-flags 1
/dev/null: [O_RDWR]
/dev/null: [O_RDWR]
/proc/1/mountinfo: [O_CLOEXEC]
/proc/swaps: [O_CLOEXEC]
/run/cloud-init/hook-hotplug-cmd: [O_RDWR O_CLOEXEC O_NDELAY]
/dev/rfkill: [O_RDWR O_CLOEXEC O_NDELAY]
/dev/null: [O_RDWR]
/dev/autofs: [O_CLOEXEC]
/dev/kmsg: [O_WRONLY O_CLOEXEC]
/sys/fs/cgroup: [O_CLOEXEC O_DIRECTORY O_NDELAY]
/run/dmeventd-server: [O_RDWR O_CLOEXEC O_NDELAY]
/run/dmeventd-client: [O_RDWR O_CLOEXEC O_NDELAY]
/run/initctl: [O_RDWR O_CLOEXEC O_NDELAY]
```

### Examining a web server process

```bash
$ show-open-flags $(pgrep nginx | head -1)
/etc/nginx/nginx.conf: [O_RDONLY O_CLOEXEC]
/var/log/nginx/access.log: [O_WRONLY O_APPEND O_CLOEXEC]
/var/log/nginx/error.log: [O_WRONLY O_APPEND O_CLOEXEC]
/usr/share/nginx/html/index.html: [O_RDONLY O_CLOEXEC]
```

## Common Flags Explained

| Flag | Description |
|------|-------------|
| `O_RDONLY` | Open for reading only |
| `O_WRONLY` | Open for writing only |
| `O_RDWR` | Open for reading and writing |
| `O_APPEND` | Append on each write |
| `O_CREATE` | Create file if it does not exist |
| `O_EXCL` | Error if CREATE and file exists |
| `O_SYNC` | Make writes synchronous |
| `O_TRUNC` | Truncate file to zero length |
| `O_ASYNC` | Enable signal-driven I/O |
| `O_CLOEXEC` | Close on exec |
| `O_DIRECT` | Direct I/O |
| `O_DIRECTORY` | Must be a directory |
| `O_DSYNC` | Synchronized I/O data integrity completion |
| `O_LARGEFILE` | Allow files whose sizes cannot be represented in an off_t |
| `O_NDELAY` | Non-blocking mode |

## How It Works

`show-open-flags` works by examining the Linux `/proc` filesystem, which provides detailed information about running processes:

1. For a given PID, it reads `/proc/[pid]/fd/` to enumerate open file descriptors
2. It resolves symlinks in `/proc/[pid]/fd/` to find the actual paths
3. It reads `/proc/[pid]/fdinfo/[fd]` to determine the open flags for each file descriptor
4. It parses the octal flags value and translates it to human-readable flag names

## Use Cases

- Debugging file access issues in applications
- Analyzing security concerns related to file access
- Understanding how processes interact with files
- Educational purposes for learning about system programming

## Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is open source. Consider adding a LICENSE file to specify terms of use.

## Acknowledgments

- Linux kernel developers for the `/proc` filesystem
- Go programming language for making system programming accessible
