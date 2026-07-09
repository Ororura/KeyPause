# KeyPause

KeyPause is a lightweight macOS menu bar app that temporarily blocks keyboard input, so you can safely clean your keyboard without accidental typing, shortcuts, or media key actions.

The app is distributed as a standard macOS `.app` bundle and can be packaged into a drag-and-drop `.dmg` installer.

## Features

- Blocks regular keyboard key down and key up events.
- Blocks modifier key changes.
- Blocks MacBook media keys such as volume and playback controls.
- Runs as a menu bar app without a Dock icon.
- Supports timed cleaning sessions from the tray menu.
- Can be stopped manually at any time.
- Also includes an optional CLI entrypoint for development or terminal usage.

## Requirements

- macOS 12 or newer
- Go 1.26.4 or newer for building from source
- Accessibility permission for `KeyPause.app`

KeyPause uses a macOS `CGEventTap`, so Accessibility permission is required. If permission is missing, macOS will prompt you to enable it in System Settings.

## Build The App

Build `KeyPause.app`:

```bash
./scripts/build-app.sh
```

The app bundle will be created at:

```text
dist/KeyPause.app
```

Run it locally:

```bash
open dist/KeyPause.app
```

## Build The DMG

Build a drag-and-drop DMG installer:

```bash
./scripts/build-dmg.sh
```

The DMG will be created at:

```text
dist/KeyPause.dmg
```

## Installation

1. Open `dist/KeyPause.dmg`.
2. Drag `KeyPause.app` to `Applications`.
3. Start `KeyPause` from `Applications`.
4. Grant Accessibility permission in System Settings when prompted.

After startup, KeyPause appears in the macOS menu bar. The tray menu includes:

- `Start for 30s`
- `Start for 1m`
- `Stop`
- `Quit`

## Code Signing

The build script uses ad-hoc signing:

```bash
codesign --force --deep --sign - dist/KeyPause.app
```

This is enough for local use. Public distribution usually requires an Apple Developer ID certificate, hardened runtime, notarization, and stapling.

## CLI Mode

Build the CLI version:

```bash
go build -o keyboard-cleaner ./cmd
```

Run a 30-second cleaning session from the terminal:

```bash
./keyboard-cleaner -time 30s
```

Run a 1-minute cleaning session:

```bash
./keyboard-cleaner -time 1m
```

## Testing

Run all tests:

```bash
go test ./...
```

## Project Structure

- `cmd/main.go` contains the CLI entrypoint.
- `cmd/tray/main.go` contains the menu bar app entrypoint.
- `build/macos/Info.plist` contains the macOS app bundle metadata.
- `scripts/build-app.sh` builds `KeyPause.app`.
- `scripts/build-dmg.sh` builds `KeyPause.dmg`.
- `internal/app` contains the application orchestration logic.
- `internal/blocker/macos` contains the native macOS event tap implementation.
- `internal/tray` contains tray controller logic and icon generation.

---

# KeyPause

KeyPause - легкое menu bar приложение для macOS, которое временно блокирует ввод с клавиатуры, чтобы можно было безопасно почистить клавиатуру без случайного набора текста, горячих клавиш или срабатывания медиа-клавиш.

Приложение собирается как стандартный macOS `.app` bundle и может быть упаковано в drag-and-drop `.dmg` установщик.

## Возможности

- Блокирует обычные нажатия и отпускания клавиш.
- Блокирует изменения modifier-клавиш.
- Блокирует медиа-клавиши MacBook, включая громкость и управление воспроизведением.
- Работает как menu bar приложение без иконки в Dock.
- Поддерживает сессии очистки по таймеру через tray-меню.
- Можно остановить блокировку вручную в любой момент.
- Также содержит optional CLI entrypoint для разработки или запуска из терминала.

## Требования

- macOS 12 или новее
- Go 1.26.4 или новее для сборки из исходников
- Accessibility permission для `KeyPause.app`

KeyPause использует macOS `CGEventTap`, поэтому приложению нужен доступ Accessibility. Если доступа нет, macOS предложит включить его в System Settings.

## Сборка Приложения

Собрать `KeyPause.app`:

```bash
./scripts/build-app.sh
```

App bundle будет создан здесь:

```text
dist/KeyPause.app
```

Запустить локально:

```bash
open dist/KeyPause.app
```

## Сборка DMG

Собрать drag-and-drop DMG установщик:

```bash
./scripts/build-dmg.sh
```

DMG будет создан здесь:

```text
dist/KeyPause.dmg
```

## Установка

1. Откройте `dist/KeyPause.dmg`.
2. Перетащите `KeyPause.app` в `Applications`.
3. Запустите `KeyPause` из `Applications`.
4. Разрешите Accessibility permission в System Settings, когда macOS попросит доступ.

После запуска KeyPause появляется в menu bar. В tray-меню доступны пункты:

- `Start for 30s`
- `Start for 1m`
- `Stop`
- `Quit`

## Code Signing

Скрипт сборки использует ad-hoc signing:

```bash
codesign --force --deep --sign - dist/KeyPause.app
```

Этого достаточно для локального использования. Для публичного распространения обычно нужны Apple Developer ID certificate, hardened runtime, notarization и stapling.

## CLI Режим

Собрать CLI-версию:

```bash
go build -o keyboard-cleaner ./cmd
```

Запустить 30-секундную сессию очистки из терминала:

```bash
./keyboard-cleaner -time 30s
```

Запустить сессию на 1 минуту:

```bash
./keyboard-cleaner -time 1m
```

## Тесты

Запустить все тесты:

```bash
go test ./...
```

## Структура Проекта

- `cmd/main.go` содержит CLI entrypoint.
- `cmd/tray/main.go` содержит entrypoint menu bar приложения.
- `build/macos/Info.plist` содержит metadata macOS app bundle.
- `scripts/build-app.sh` собирает `KeyPause.app`.
- `scripts/build-dmg.sh` собирает `KeyPause.dmg`.
- `internal/app` содержит orchestration-логику приложения.
- `internal/blocker/macos` содержит нативную macOS event tap реализацию.
- `internal/tray` содержит controller tray-режима и генерацию иконки.
