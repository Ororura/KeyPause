package macos

import (
	"context"
	"errors"
	"sync"
)

type KeyboardBlocker struct {
	mu     sync.Mutex
	active bool
}

func NewKeyboardBlocker() *KeyboardBlocker {
	return &KeyboardBlocker{}
}

func (b *KeyboardBlocker) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.active {
		return nil
	}

	if !accessibilityTrusted(false) {
		accessibilityTrusted(true)
		return errors.New("accessibility permission is required; enable it for this terminal app and keyboard-cleaner in System Settings")
	}

	ok := startEventTap()
	if !ok {
		return errors.New("failed to start keyboard blocker: check Accessibility permissions")
	}

	b.active = true

	go func() {
		<-ctx.Done()
		_ = b.Stop()
	}()

	return nil
}

func (b *KeyboardBlocker) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	stopEventTap()
	b.active = false

	return nil
}

func (b *KeyboardBlocker) IsActive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.active
}
