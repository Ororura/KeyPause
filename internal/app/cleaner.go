package app

import (
	"context"
	"time"
)

type KeyboardBlocker interface {
	Start(ctx context.Context) error
	Stop() error
	IsActive() bool
}

type Cleaner struct {
	blocker KeyboardBlocker
}

func NewCleaner(blocker KeyboardBlocker) *Cleaner {
	return &Cleaner{blocker: blocker}
}

func (c *Cleaner) StartCleaning(ctx context.Context, duration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	if err := c.blocker.Start(ctx); err != nil {
		return err
	}

	<-ctx.Done()

	return c.blocker.Stop()
}
