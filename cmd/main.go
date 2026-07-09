package main

import (
	"context"
	"flag"
	"fmt"
	"go-keyboard-cleaner/internal/app"
	"go-keyboard-cleaner/internal/blocker/macos"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	duration := flag.Duration("time", 30*time.Second, "cleaning duration, example: 30s, 1m")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	keyboardBlocker := macos.NewKeyboardBlocker()
	cleaner := app.NewCleaner(keyboardBlocker)

	fmt.Printf("Keyboard cleaning started for %s\n", duration.String())
	fmt.Println("Press Ctrl+C to sleep earlier.")

	if err := cleaner.StartCleaning(ctx, *duration); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Keyboard cleaning stopped.")
}
