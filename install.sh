#!/bin/bash
set -e

REPO="victorrgr/battery-monitor"
BINARY_NAME="battery-monitor"
INSTALL_DIR="$HOME/.local/bin"
AUTOSTART_DIR="$HOME/.config/autostart"
DESKTOP_FILE="$AUTOSTART_DIR/$BINARY_NAME.desktop"
DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME"

echo "🔧 Creating install directories..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$AUTOSTART_DIR"

echo "⬇️ Downloading $BINARY_NAME from GitHub Releases..."
curl -sSfL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "🖥 Creating autostart desktop entry..."
cat > "$DESKTOP_FILE" <<EOF
[Desktop Entry]
Type=Application
Exec=$INSTALL_DIR/$BINARY_NAME
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name=Battery Monitor
Comment=Monitors battery usage in background
EOF

echo "✅ Installed to: $INSTALL_DIR/$BINARY_NAME"
echo "✅ Autostart file: $DESKTOP_FILE"
echo "📈 You can now run \`$BINARY_NAME analyse\` from your terminal."