#!/bin/bash

# Configuration
PUBLIC_URL=$1
APP_DIR="src/frontend/teams-app"
PUBLIC_DIR="$APP_DIR/public"
OUTPUT_FILE="app.zip"

if [ -z "$PUBLIC_URL" ]; then
    echo "Usage: sh package-app.sh <https-tunnel-url>"
    echo "Example: sh package-app.sh https://random-id.ngrok-free.app"
    exit 1
fi

# Clean up protocol for valid domains
DOMAIN=$(echo $PUBLIC_URL | sed -e 's|^https://||' -e 's|/.*$||')

echo "Packaging Teams App..."
echo "Public URL: $PUBLIC_URL"
echo "Domain: $DOMAIN"

# 1. Update manifest.json placeholders (temp file)
TEMP_MANIFEST=$(mktemp)
sed -e "s|https://YOUR_NGROK_URL|$PUBLIC_URL|g" \
    -e "s|YOUR_NGROK_DOMAIN|$DOMAIN|g" \
    "$PUBLIC_DIR/manifest.json" > "$TEMP_MANIFEST"

# 2. Check for icons
if [ ! -f "$PUBLIC_DIR/color.png" ] || [ ! -f "$PUBLIC_DIR/outline.png" ]; then
    echo "Error: Required icons (color.png, outline.png) not found in $PUBLIC_DIR"
    exit 1
fi

# 3. Create zip
# We need to stay in the directory or use relative paths correctly
# Manifest and icons must be at the root of the zip
# First, create a clean build directory for the package
mkdir -p build/teams-pkg
cp "$TEMP_MANIFEST" build/teams-pkg/manifest.json
cp "$PUBLIC_DIR/color.png" build/teams-pkg/
cp "$PUBLIC_DIR/outline.png" build/teams-pkg/

# Create the zip in the root directory
cd build/teams-pkg
zip -r "../../app.zip" *
cd - > /dev/null

# Clean up
rm -rf build/teams-pkg
rm "$TEMP_MANIFEST"

echo "Success! Created app.zip in the project root."
echo "You can now upload this file to Microsoft Teams."

