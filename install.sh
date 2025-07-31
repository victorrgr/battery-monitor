#!/bin/bash
set -e

REPO="victorrgr/battery-monitor"
BINARY_NAME="battery-monitor"
INSTALL_DIR="$HOME/.local/bin"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
SERVICE_FILE="$SYSTEMD_USER_DIR/$BINARY_NAME.service"
DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME"

echo "🔧 Creating install directories..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$SYSTEMD_USER_DIR"

echo "⬇️ Downloading $BINARY_NAME from GitHub Releases..."
curl -sSfL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "🔍 Ensuring $INSTALL_DIR is in PATH..."
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    SHELL_RC=""
    if [ -n "$BASH_VERSION" ]; then
        SHELL_RC="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="$HOME/.zshrc"
    else
        SHELL_RC="$HOME/.profile"
    fi

    if ! grep -q "$INSTALL_DIR" "$SHELL_RC"; then
        echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$SHELL_RC"
        echo "✅ Added $INSTALL_DIR to PATH in $SHELL_RC"
    fi
else
    echo "✅ $INSTALL_DIR is already in PATH"
fi

echo "🛠 Creating systemd user service file..."
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Battery Monitor

[Service]
ExecStart=$INSTALL_DIR/$BINARY_NAME monitor
Restart=on-failure

[Install]
WantedBy=default.target
EOF

echo "🔄 Reloading systemd user units..."
systemctl --user daemon-reexec
systemctl --user daemon-reload

echo "✅ Enabling and starting battery-monitor.service..."
systemctl --user enable --now "$BINARY_NAME.service"

echo
echo "🎉 Installation complete!"
echo "📦 Installed to: $INSTALL_DIR/$BINARY_NAME"
echo "🛠 Service created: $SERVICE_FILE"
echo "🚀 Monitoring will start automatically on login."
echo "📈 Run \`$BINARY_NAME analyse\` to generate a battery report."
echo "🔍 Use \`systemctl --user status $BINARY_NAME.service\` to check monitor status."