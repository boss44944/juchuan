#!/bin/bash
set -e

APP_NAME="Juchuan"
BINARY="juchuan"
APP_DIR="${APP_NAME}.app"
CONTENTS="${APP_DIR}/Contents"

rm -rf "${APP_DIR}"
mkdir -p "${CONTENTS}/MacOS" "${CONTENTS}/Resources"

GOOS=darwin GOARCH=arm64 go build -o "${CONTENTS}/MacOS/${BINARY}" .

if [ -f "app-logo.png" ]; then
  mkdir -p icon.iconset
  sips -z 16 16 app-logo.png --out icon.iconset/icon_16x16.png >/dev/null
  sips -z 32 32 app-logo.png --out icon.iconset/icon_16x16@2x.png >/dev/null
  sips -z 32 32 app-logo.png --out icon.iconset/icon_32x32.png >/dev/null
  sips -z 64 64 app-logo.png --out icon.iconset/icon_32x32@2x.png >/dev/null
  sips -z 128 128 app-logo.png --out icon.iconset/icon_128x128.png >/dev/null
  sips -z 256 256 app-logo.png --out icon.iconset/icon_128x128@2x.png >/dev/null
  sips -z 256 256 app-logo.png --out icon.iconset/icon_256x256.png >/dev/null
  sips -z 512 512 app-logo.png --out icon.iconset/icon_256x256@2x.png >/dev/null
  sips -z 512 512 app-logo.png --out icon.iconset/icon_512x512.png >/dev/null
  sips -z 1024 1024 app-logo.png --out icon.iconset/icon_512x512@2x.png >/dev/null
  iconutil -c icns icon.iconset -o "${CONTENTS}/Resources/Juchuan.icns"
  rm -rf icon.iconset
fi

cat > "${CONTENTS}/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>Juchuan</string>
    <key>CFBundleDisplayName</key>
    <string>菊传</string>
    <key>CFBundleIdentifier</key>
    <string>com.juchuan.app</string>
    <key>CFBundleExecutable</key>
    <string>${BINARY}</string>
    <key>CFBundleIconFile</key>
    <string>Juchuan.icns</string>
</dict>
</plist>
EOF

chmod +x "${CONTENTS}/MacOS/${BINARY}"

echo "Created ${APP_DIR}"