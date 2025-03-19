# pings - Multiple Ping Monitoring Tool

## Overview
`pings` is a simple command-line tool that continuously pings a list of IP addresses and shows their online/offline status. It also tracks successful and failed pings over time.

## Features
- Ping multiple IP addresses at once
- Read IPs from a file
- Adjustable scan interval
- Color-coded output for online/offline status
- Tracks successful (`ok`) and failed (`!ok`) pings

## Installation
Make sure you have **Go** installed and run:

```sh
go build -o pings pings.go
```

This will create an executable `pings` (or `pings.exe` on Windows).

## Usage

```sh
pings [IP1 IP2 ...] [interval]
pings [file] [interval]
```

### Examples:
- **Ping multiple IPs with default interval (5s):**
  ```sh
  pings 192.168.1.1 192.168.1.2
  ```
- **Ping multiple IPs with a 10s interval:**
  ```sh
  pings 192.168.1.1 192.168.1.2 10
  ```
- **Read IPs from a file and use a 5s interval:**
  ```sh
  pings ip-list.txt
  ```
- **Read IPs from a file and use a 15s interval:**
  ```sh
  pings ip-list.txt 15
  ```

## Output Example
```
--- IP Status Check ---
192.168.1.15:     Offline   (0 ok / 2 !ok)
192.168.1.254:    Online    (5 ok / 0 !ok)
192.168.1.40:     Online    (3 ok / 0 !ok)
```
- **Green (✔ Online)** = IP is reachable
- **Red (✖ Offline)** = IP is unreachable
- **Yellow counter** = The IP has at least one failed ping

## Notes
- On Windows, ANSI color codes may not work in `cmd.exe`. Use **PowerShell** or **Windows Terminal**.
- Requires Go to compile. Precompiled binaries can be created for Linux, Windows, and macOS using:
  ```sh
  GOOS=linux GOARCH=amd64 go build -o pings-linux pings.go
  GOOS=windows GOARCH=amd64 go build -o pings.exe pings.go
  GOOS=darwin GOARCH=amd64 go build -o pings-mac pings.go
  ```

