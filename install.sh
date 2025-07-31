#/bin/bash
set -e

REPO=victorrgr/battery-monitor
BINARY_NAME=battery-monitor
INSTALL_DIR=/home/victorrgr/.local/bin
AUTOSTART_DIR=/home/victorrgr/.config/autostart
DESKTOP_FILE=/.desktop
DOWNLOAD_URL=https://github.com//releases/latest/download/

echo 🔧 Creating install directories...
mkdir -p 
mkdir -p 

echo ⬇️ Downloading from GitHub Releases...
curl -sSfL  -o /
chmod +x /

echo 🖥 Creating autostart desktop entry...
cat >  <<EOF
[Desktop Entry]
Type=Application
Exec=/
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name=Battery Monitor
Comment=Monitors battery usage in background
EOF

echo ✅ Installed to: /
echo ✅ Autostart file: 
echo 📈 You can now run ` analyse` from your terminal.

