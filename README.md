# Battery Monitor

A lightweight CLI tool that records battery usage into a local SQLite database and generates reports in HTML format.

## Features

- Runs in the background after login
- Logs battery percentage and status over time
- Generates reports for historical analysis

## Installation

For Ubuntu or any systemd-based system:

```bash
curl -sSfL https://raw.githubusercontent.com/victorrgr/battery-monitor/master/install.sh | bash
```

This installs the binary to `~/.local/bin`, creates a systemd user service, and enables it.

## Usage

```bash
battery-monitor monitor     # Start monitoring manually (optional)
battery-monitor analyse     # Generate report.html in the current directory
xdg-open report.html        # Open the report in a browser
```

Check service status:

```bash
systemctl --user status battery-monitor.service
```

## Uninstall

### Current version (systemd):

```bash
curl -sSfL https://raw.githubusercontent.com/victorrgr/battery-monitor/master/uninstall.sh | bash
```

### Legacy version (.desktop file):

```bash
curl -sSfL https://raw.githubusercontent.com/victorrgr/battery-monitor/master/uninstall-legacy.sh | bash
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

## Data and Reports

- Data: `~/.local/share/battery-monitor/battery-monitor.db`
- Reports: `./report.html`
