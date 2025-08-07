# Battery Monitor

A lightweight CLI tool that records battery usage into a local SQLite database and exposes a web interface for viewing battery data.

## Features

- Runs in the background after login
- Logs battery percentage and status over time
- Serves a local web dashboard for historical analysis

## Installation

For Ubuntu or any systemd-based system:

```bash
curl -sSfL https://raw.githubusercontent.com/victorrgr/battery-monitor/master/install.sh | bash
```

This installs the binary to `~/.local/bin`, creates a systemd user service, and enables it.

## Usage

```bash
battery-monitor monitor     # Start monitoring manually (optional)
battery-monitor analyse     # Start the web server (default port: 8080)
```

> To view the battery data UI, you must run the `analyse` command in your terminal.  
> Once it's running, open [http://localhost:8080](http://localhost:8080) in your browser.

Check service status:

```bash
systemctl --user status battery-monitor.service
```

## Uninstall

```bash
curl -sSfL https://raw.githubusercontent.com/victorrgr/battery-monitor/master/uninstall.sh | bash
```

## Build from Source

### Requirements

- Go 1.22 or newer
- SQLite development headers:
  ```bash
  sudo apt install build-essential libsqlite3-dev
  ```

### Build

```bash
CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/battery-monitor ./cmd/battery-monitor
```

## Data Location

- Battery data is stored in: `~/.local/share/battery-monitor/battery-monitor.db`