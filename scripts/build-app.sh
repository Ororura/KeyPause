#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_NAME="KeyPause"
DIST_DIR="${ROOT_DIR}/dist"
APP_DIR="${DIST_DIR}/${APP_NAME}.app"
CONTENTS_DIR="${APP_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"
ICONSET_DIR="${DIST_DIR}/${APP_NAME}.iconset"

rm -rf "${APP_DIR}" "${ICONSET_DIR}"
mkdir -p "${MACOS_DIR}" "${RESOURCES_DIR}" "${ICONSET_DIR}"

go build -o "${MACOS_DIR}/${APP_NAME}" "${ROOT_DIR}/cmd/tray"
cp "${ROOT_DIR}/build/macos/Info.plist" "${CONTENTS_DIR}/Info.plist"

ICON_SRC="${DIST_DIR}/${APP_NAME}-1024.png"
ICON_GEN="${DIST_DIR}/generate-icon.go"
cat > "${ICON_GEN}" <<'GOEOF'
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	img := image.NewNRGBA(image.Rect(0, 0, 1024, 1024))
	bg := color.NRGBA{R: 106, G: 90, B: 255, A: 255}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	for y := 0; y < 1024; y++ {
		for x := 0; x < 1024; x++ {
			img.SetNRGBA(x, y, bg)
		}
	}
	fill(img, 180, 290, 844, 640, white)
	fill(img, 245, 355, 779, 575, bg)
	fill(img, 290, 690, 734, 790, white)
	fill(img, 430, 790, 594, 850, white)
	fill(img, 315, 395, 410, 485, white)
	fill(img, 465, 395, 560, 485, white)
	fill(img, 615, 395, 710, 485, white)
	fill(img, 315, 525, 710, 575, white)

	file, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func fill(img *image.NRGBA, x0, y0, x1, y1 int, c color.NRGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.SetNRGBA(x, y, c)
		}
	}
}
GOEOF

go run "${ICON_GEN}" "${ICON_SRC}"
rm "${ICON_GEN}"

for size in 16 32 128 256 512; do
	sips -z "${size}" "${size}" "${ICON_SRC}" --out "${ICONSET_DIR}/icon_${size}x${size}.png" >/dev/null
	double_size=$((size * 2))
	sips -z "${double_size}" "${double_size}" "${ICON_SRC}" --out "${ICONSET_DIR}/icon_${size}x${size}@2x.png" >/dev/null
done

iconutil -c icns "${ICONSET_DIR}" -o "${RESOURCES_DIR}/${APP_NAME}.icns"
rm -rf "${ICONSET_DIR}" "${ICON_SRC}"

codesign --force --deep --sign - "${APP_DIR}" >/dev/null

echo "Built ${APP_DIR}"
