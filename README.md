# KeyPause

KeyPause is a lightweight macOS utility that temporarily blocks keyboard input, so you can safely clean your keyboard without accidental typing, shortcuts, or media key actions.

It can run as a command-line tool for timed cleaning sessions or as a menu bar app with quick start/stop controls.

## Features

- Blocks regular keyboard key down and key up events.
- Blocks modifier key changes.
- Blocks MacBook media keys such as volume and playback controls.
- Supports timed cleaning sessions.
- Can be stopped early with `Ctrl+C` in CLI mode or from the tray menu.
- Uses a macOS menu bar tray app for quick access.

## Requirements

- macOS
- Go 1.26.4 or newer
- Accessibility permission for the terminal app and/or the built binary

KeyPause uses a macOS `CGEventTap`, so Accessibility permission is required. If permission is missing, macOS will prompt you to enable it in System Settings.

## Build

Build the CLI version:

```bash
go build -o keyboard-cleaner ./cmd
```

Build the menu bar version:

```bash
go build -o keyboard-cleaner-tray ./cmd/tray
```

## Usage

Run a 30-second cleaning session from the terminal:

```bash
./keyboard-cleaner -time 30s
```

Run a 1-minute cleaning session:

```bash
./keyboard-cleaner -time 1m
```

Start the menu bar app:

```bash
./keyboard-cleaner-tray
```

The tray menu includes:

- `Start for 30s`
- `Start for 1m`
- `Stop`
- `Quit`

## Testing

Run all tests:

```bash
go test ./...
```

## Project Structure

- `cmd/main.go` contains the CLI entrypoint.
- `cmd/tray/main.go` contains the menu bar app entrypoint.
- `internal/app` contains the application orchestration logic.
- `internal/blocker/macos` contains the native macOS event tap implementation.
- `internal/tray` contains tray controller logic and icon generation.

---

# KeyPause

KeyPause - легкая утилита для macOS, которая временно блокирует ввод с клавиатуры, чтобы можно было безопасно почистить клавиатуру без случайного набора текста, горячих клавиш или срабатывания медиа-клавиш.

Приложение можно запускать как CLI-утилиту с таймером или как приложение в menu bar с быстрым управлением через tray-меню.

## Возможности

- Блокирует обычные нажатия и отпускания клавиш.
- Блокирует изменения modifier-клавиш.
- Блокирует медиа-клавиши MacBook, включая громкость и управление воспроизведением.
- Поддерживает сессии очистки по таймеру.
- Можно остановить раньше через `Ctrl+C` в CLI-режиме или через tray-меню.
- Есть menu bar приложение для быстрого запуска.

## Требования

- macOS
- Go 1.26.4 или новее
- Accessibility permission для терминала и/или собранного бинарника

KeyPause использует macOS `CGEventTap`, поэтому приложению нужен доступ Accessibility. Если доступа нет, macOS предложит включить его в System Settings.

## Сборка

Собрать CLI-версию:

```bash
go build -o keyboard-cleaner ./cmd
```

Собрать menu bar версию:

```bash
go build -o keyboard-cleaner-tray ./cmd/tray
```

## Использование

Запустить 30-секундную сессию очистки из терминала:

```bash
./keyboard-cleaner -time 30s
```

Запустить сессию на 1 минуту:

```bash
./keyboard-cleaner -time 1m
```

Запустить menu bar приложение:

```bash
./keyboard-cleaner-tray
```

В tray-меню доступны пункты:

- `Start for 30s`
- `Start for 1m`
- `Stop`
- `Quit`

## Тесты

Запустить все тесты:

```bash
go test ./...
```

## Структура проекта

- `cmd/main.go` содержит CLI entrypoint.
- `cmd/tray/main.go` содержит entrypoint menu bar приложения.
- `internal/app` содержит orchestration-логику приложения.
- `internal/blocker/macos` содержит нативную macOS event tap реализацию.
- `internal/tray` содержит controller tray-режима и генерацию иконки.
