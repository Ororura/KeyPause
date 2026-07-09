package app

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestStartCleaningStopsAfterDuration(t *testing.T) {
	blocker := &fakeKeyboardBlocker{}
	cleaner := NewCleaner(blocker)

	err := cleaner.StartCleaning(context.Background(), time.Millisecond)
	if err != nil {
		t.Fatalf("StartCleaning() error = %v, want nil", err)
	}

	if blocker.startCalls != 1 {
		t.Fatalf("Start() calls = %d, want 1", blocker.startCalls)
	}
	if blocker.stopCalls != 1 {
		t.Fatalf("Stop() calls = %d, want 1", blocker.stopCalls)
	}
	if blocker.active {
		t.Fatal("blocker is active after StartCleaning returned")
	}
}

func TestStartCleaningStopsWhenContextCancelled(t *testing.T) {
	blocker := &fakeKeyboardBlocker{started: make(chan struct{})}
	cleaner := NewCleaner(blocker)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- cleaner.StartCleaning(ctx, time.Hour)
	}()

	select {
	case <-blocker.started:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Start()")
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("StartCleaning() error = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("StartCleaning() did not return after context cancellation")
	}

	if blocker.stopCalls != 1 {
		t.Fatalf("Stop() calls = %d, want 1", blocker.stopCalls)
	}
}

func TestStartCleaningReturnsStartError(t *testing.T) {
	startErr := errors.New("start failed")
	blocker := &fakeKeyboardBlocker{startErr: startErr}
	cleaner := NewCleaner(blocker)

	err := cleaner.StartCleaning(context.Background(), time.Hour)
	if !errors.Is(err, startErr) {
		t.Fatalf("StartCleaning() error = %v, want %v", err, startErr)
	}

	if blocker.stopCalls != 0 {
		t.Fatalf("Stop() calls = %d, want 0", blocker.stopCalls)
	}
}

func TestStartCleaningReturnsStopError(t *testing.T) {
	stopErr := errors.New("stop failed")
	blocker := &fakeKeyboardBlocker{stopErr: stopErr}
	cleaner := NewCleaner(blocker)

	err := cleaner.StartCleaning(context.Background(), time.Millisecond)
	if !errors.Is(err, stopErr) {
		t.Fatalf("StartCleaning() error = %v, want %v", err, stopErr)
	}
}

type fakeKeyboardBlocker struct {
	startErr error
	stopErr  error
	active   bool

	startCalls int
	stopCalls  int
	started    chan struct{}
}

func (b *fakeKeyboardBlocker) Start(context.Context) error {
	b.startCalls++
	if b.started == nil {
		b.started = make(chan struct{})
	}
	close(b.started)

	if b.startErr != nil {
		return b.startErr
	}

	b.active = true
	return nil
}

func (b *fakeKeyboardBlocker) Stop() error {
	b.stopCalls++
	b.active = false
	return b.stopErr
}

func (b *fakeKeyboardBlocker) IsActive() bool {
	return b.active
}
