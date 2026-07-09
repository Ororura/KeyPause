#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_NAME="KeyPause"
DIST_DIR="${ROOT_DIR}/dist"
APP_DIR="${DIST_DIR}/${APP_NAME}.app"
DMG_PATH="${DIST_DIR}/${APP_NAME}.dmg"
TEMP_DMG="${DIST_DIR}/${APP_NAME}-temp.dmg"
VOLUME_NAME="${APP_NAME}"

"${ROOT_DIR}/scripts/build-app.sh"

rm -f "${DMG_PATH}" "${TEMP_DMG}"
hdiutil create -volname "${VOLUME_NAME}" -srcfolder "${APP_DIR}" -ov -format UDRW "${TEMP_DMG}" >/dev/null

DEVICE="$(hdiutil attach -readwrite -noverify -noautoopen "${TEMP_DMG}" | awk '/\/Volumes\// {print $1; exit}')"
VOLUME_PATH="/Volumes/${VOLUME_NAME}"

ln -s /Applications "${VOLUME_PATH}/Applications"
hdiutil detach "${DEVICE}" >/dev/null
hdiutil convert "${TEMP_DMG}" -format UDZO -imagekey zlib-level=9 -o "${DMG_PATH}" >/dev/null
rm -f "${TEMP_DMG}"

echo "Built ${DMG_PATH}"
