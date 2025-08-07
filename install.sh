#!/bin/bash
set -e

REPO="victorrgr/battery-monitor"
BINARY_NAME="battery-monitor"
INSTALL_DIR="$HOME/.local/bin"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
SERVICE_FILE="$SYSTEMD_USER_DIR/$BINARY_NAME.service"
DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME"

FORCE_UPDATE=false
if [[ "$1" == "--force" || "$1" == "--update" ]]; then
    FORCE_UPDATE=true
fi

echo "🔧 Creating install directories..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$SYSTEMD_USER_DIR"

if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    if [ "$FORCE_UPDATE" = true ]; then
        echo "⬇️ Updating $BINARY_NAME from GitHub Releases..."
        if systemctl --user is-active --quiet "$BINARY_NAME.service"; then
            echo "⏸ Stopping running service before update..."
            systemctl --user stop "$BINARY_NAME.service"
        fi

        rm -f "$INSTALL_DIR/$BINARY_NAME"
        curl -sSfL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
        echo "✅ Updated $BINARY_NAME"

        if systemctl --user is-enabled --quiet "$BINARY_NAME.service"; then
            echo "🚀 Restarting service after update..."
            systemctl --user start "$BINARY_NAME.service"
            echo "✅ Service restarted"
        fi
    else
        echo "✅ $BINARY_NAME is already installed. Use --update to force update."
    fi
    exit 0
else
    echo "⬇️ Downloading $BINARY_NAME from GitHub Releases..."
    curl -sSfL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    echo "✅ Installed $BINARY_NAME"
fi

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

if [ ! -f "$SERVICE_FILE" ]; then
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
else
    echo "✅ Systemd service already exists."
fi

echo
echo "🎉 Installation complete!"
echo "📦 Installed to: $INSTALL_DIR/$BINARY_NAME"
echo "🛠 Service file: $SERVICE_FILE"
echo "🚀 Monitoring will start automatically on login."
echo "📈 Run \`$BINARY_NAME analyse\` to open the web dashboard."
echo "🔍 Use \`systemctl --user status $BINARY_NAME.service\` to check monitor status."