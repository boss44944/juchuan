#!/bin/bash
set -euo pipefail

if [[ "$(uname -s)" != "Darwin" || "$(uname -m)" != "arm64" ]]; then
  echo "Juchuan.app packaging requires a native Apple Silicon macOS environment." >&2
  exit 1
fi

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
REPO_ROOT=$(cd -- "${SCRIPT_DIR}/.." && pwd -P)
DESTINATION="$(pwd -P)/Juchuan.app"
if [[ -e "${DESTINATION}" || -L "${DESTINATION}" ]]; then
  echo "Move the existing Juchuan.app aside before packaging." >&2
  exit 1
fi

STAGE=$(mktemp -d "$(pwd -P)/.juchuan-app.XXXXXX")
trap 'rm -rf -- "${STAGE}"' EXIT
trap 'exit 1' HUP INT TERM

BINARY="juchuan"
APP_DIR="${STAGE}/Juchuan.app"
CONTENTS="${APP_DIR}/Contents"
ICON_SOURCE="${SCRIPT_DIR}/app-logo.png"
ICON_SET="${STAGE}/Juchuan.iconset"
mkdir -p "${CONTENTS}/MacOS" "${CONTENTS}/Resources" "${ICON_SET}"

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 MACOSX_DEPLOYMENT_TARGET=13.0 \
  "${REPO_ROOT}/scripts/build-single-binary.sh" "${CONTENTS}/MacOS/${BINARY}"

for size in 16 32 128 256 512; do
  sips -z "${size}" "${size}" "${ICON_SOURCE}" \
    --out "${ICON_SET}/icon_${size}x${size}.png" >/dev/null
  sips -z "$((size * 2))" "$((size * 2))" "${ICON_SOURCE}" \
    --out "${ICON_SET}/icon_${size}x${size}@2x.png" >/dev/null
done
iconutil -c icns "${ICON_SET}" -o "${CONTENTS}/Resources/Juchuan.icns"

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
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>JCHN</string>
    <key>CFBundleSupportedPlatforms</key>
    <array>
        <string>MacOSX</string>
    </array>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>13.0</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

plutil -lint "${CONTENTS}/Info.plist"
test -x "${CONTENTS}/MacOS/${BINARY}"
[[ "$(lipo -archs "${CONTENTS}/MacOS/${BINARY}")" == "arm64" ]]

# Ad hoc signing is enough for verifying a local build artifact, not public distribution.
codesign --force --sign - "${APP_DIR}"
codesign --verify --deep --strict "${APP_DIR}"
mv "${APP_DIR}" "${DESTINATION}"
echo "Created ${DESTINATION}"
