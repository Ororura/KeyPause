package tray

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestControllerStartStopsAfterDuration(t *testing.T) {
	blocker := &fakeBlocker{}
	controller := NewController(blocker, nil)

	if err := controller.Start(time.Millisecond); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	if !controller.IsActive() {
		t.Fatal("controller is inactive after Start()")
	}

	eventually(t, func() bool { return !controller.IsActive() })

	if got := blocker.startCallCount(); got != 1 {
		t.Fatalf("Start() calls = %d, want 1", got)
	}
	if got := blocker.stopCallCount(); got != 1 {
		t.Fatalf("Stop() calls = %d, want 1", got)
	}
}

func TestControllerStopStopsActiveRun(t *testing.T) {
	blocker := &fakeBlocker{}
	controller := NewController(blocker, nil)

	if err := controller.Start(time.Hour); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	if err := controller.Stop(); err != nil {
		t.Fatalf("Stop() error = %v, want nil", err)
	}

	if controller.IsActive() {
		t.Fatal("controller is active after Stop()")
	}
	if got := blocker.stopCallCount(); got != 1 {
		t.Fatalf("Stop() calls = %d, want 1", got)
	}
}

func TestControllerStartReturnsStartError(t *testing.T) {
	startErr := errors.New("start failed")
	blocker := &fakeBlocker{startErr: startErr}
	controller := NewController(blocker, nil)

	err := controller.Start(time.Hour)
	if !errors.Is(err, startErr) {
		t.Fatalf("Start() error = %v, want %v", err, startErr)
	}
	if controller.IsActive() {
		t.Fatal("controller is active after failed Start()")
	}
	if got := blocker.stopCallCount(); got != 0 {
		t.Fatalf("Stop() calls = %d, want 0", got)
	}
}

func eventually(t *testing.T, condition func() bool) {
	t.Helper()

	deadline := time.After(time.Second)
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-deadline:
			t.Fatal("condition was not met before timeout")
		case <-tick.C:
			if condition() {
				return
			}
		}
	}
}

type fakeBlocker struct {
	mu sync.Mutex

	startErr error
	active   bool

	startCalls int
	stopCalls  int
}

func (b *fakeBlocker) Start(context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.startCalls++
	if b.startErr != nil {
		return b.startErr
	}

	b.active = true
	return nil
}

func (b *fakeBlocker) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.stopCalls++
	b.active = false
	return nil
}

func (b *fakeBlocker) IsActive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.active
}

func (b *fakeBlocker) startCallCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.startCalls
}

func (b *fakeBlocker) stopCallCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.stopCalls
}
