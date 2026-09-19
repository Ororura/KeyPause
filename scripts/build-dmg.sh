#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_NAME="KeyPause"
DIST_DIR="${ROOT_DIR}/dist"
APP_DIR="${DIST_DIR}/${APP_NAME}.app"
DMG_PATH="${DIST_DIR}/${APP_NAME}.dmg"
VOLUME_NAME="${APP_NAME}"

"${ROOT_DIR}/scripts/build-app.sh"

# Build the drag-and-drop DMG from a staging directory instead of mounting
# a writable image. GitHub-hosted macOS runners may hold mounted images open,
# causing hdiutil detach to fail with "Resource busy".
STAGING_DIR="$(mktemp -d "${TMPDIR:-/tmp}/keypause-dmg.XXXXXX")"
trap 'rm -rf "$STAGING_DIR"' EXIT

ditto "${APP_DIR}" "${STAGING_DIR}/${APP_NAME}.app"
ln -s /Applications "${STAGING_DIR}/Applications"

rm -f "${DMG_PATH}"
hdiutil create \
  -volname "${VOLUME_NAME}" \
  -srcfolder "${STAGING_DIR}" \
  -ov \
  -format UDZO \
  "${DMG_PATH}" >/dev/null

echo "Built ${DMG_PATH}"
