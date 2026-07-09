package tray

import (
	"context"
	"fmt"
	"go-keyboard-cleaner/internal/app"
	"sync"
	"time"
)

type State struct {
	Active bool
	Status string
	Error  string
}

type Controller struct {
	blocker app.KeyboardBlocker
	onState func(State)

	mu     sync.Mutex
	active bool
	cancel context.CancelFunc
	runID  int
}

func NewController(blocker app.KeyboardBlocker, onState func(State)) *Controller {
	return &Controller{blocker: blocker, onState: onState}
}

func (c *Controller) Start(duration time.Duration) error {
	if err := c.Stop(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	if err := c.blocker.Start(context.Background()); err != nil {
		cancel()
		c.publish(State{Status: "Stopped", Error: err.Error()})
		return err
	}

	c.mu.Lock()
	c.active = true
	c.cancel = cancel
	c.runID++
	runID := c.runID
	c.mu.Unlock()

	c.publish(State{Active: true, Status: fmt.Sprintf("Cleaning for %s", duration)})

	go func() {
		<-ctx.Done()
		c.stopRun(runID)
	}()

	return nil
}

func (c *Controller) Stop() error {
	return c.stopRun(0)
}

func (c *Controller) IsActive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.active
}

func (c *Controller) stopRun(runID int) error {
	c.mu.Lock()
	if runID != 0 && runID != c.runID {
		c.mu.Unlock()
		return nil
	}
	if !c.active {
		c.mu.Unlock()
		return nil
	}

	c.active = false
	c.runID++
	cancel := c.cancel
	c.cancel = nil
	c.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	if err := c.blocker.Stop(); err != nil {
		c.publish(State{Status: "Stopped", Error: err.Error()})
		return err
	}

	c.publish(State{Status: "Stopped"})
	return nil
}

func (c *Controller) publish(state State) {
	if c.onState != nil {
		c.onState(state)
	}
}
