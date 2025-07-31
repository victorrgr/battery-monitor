# Battery Monitor
Records battery usage to a SQLite database so it can be analysed after.

## Install
# Requirements
- GO version `1.25`

## Ubuntu
`curl -sSfL https://raw.githubusercontent.com/victorrgr/battery-monitor/main/install.sh | bash
`

## Build
### Requirements
- GO version `1.25`
- `sudo apt install build-essential libsqlite3-dev`

Command to create the build:
```bash
CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/battery-monitor ./cmd/battery-monitor
```

## Usage 
- `monitor` starts recording battery usage.
- `analyse` generates an HTML report to visualize the battery usage.