#!/bin/sh
set -e

# Default app if APP_NAME not specified
DEFAULT_APP="firstapp"

# Generate path
APP_NAME="${APP_NAME:-$DEFAULT_APP}"
APP_PATH="/app/bin/${APP_NAME}"

# Run the specified binary
exec "${APP_PATH}" "$@"