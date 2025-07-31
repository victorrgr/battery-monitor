#!/bin/bash
set -e

BINARY_NAME="battery-monitor"
INSTALL_DIR="$HOME/.local/bin"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
SERVICE_FILE="$SYSTEMD_USER_DIR/$BINARY_NAME.service"

echo "🛑 Stopping and disabling systemd user service..."
systemctl --user stop "$BINARY_NAME.service" 2>/dev/null || true
systemctl --user disable "$BINARY_NAME.service" 2>/dev/null || true

echo "🧹 Removing service file..."
rm -f "$SERVICE_FILE"

echo "🗑 Removing binary from $INSTALL_DIR..."
rm -f "$INSTALL_DIR/$BINARY_NAME"

echo "ℹ️ If you added ~/.local/bin to your PATH manually, you can remove the line from your shell config (e.g. ~/.bashrc or ~/.zshrc), but it's safe to leave it."

echo
echo "✅ Uninstall complete."
