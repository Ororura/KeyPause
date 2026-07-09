package main

import (
	"fmt"
	"go-keyboard-cleaner/internal/blocker/macos"
	"go-keyboard-cleaner/internal/tray"
	"log"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	log.SetFlags(0)
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(tray.IconPNG(), tray.IconPNG())
	systray.SetTooltip("Keyboard Cleaner")

	status := systray.AddMenuItem("Stopped", "Current keyboard cleaner status")
	status.Disable()
	systray.AddSeparator()

	start30s := systray.AddMenuItem("Start for 30s", "Block keyboard input for 30 seconds")
	start1m := systray.AddMenuItem("Start for 1m", "Block keyboard input for 1 minute")
	start15m := systray.AddMenuItem("start for 15m", "Block keyboard input for 15 minute")
	stop := systray.AddMenuItem("Stop", "Stop blocking keyboard input")
	stop.Disable()
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Quit Keyboard Cleaner")

	controller := tray.NewController(macos.NewKeyboardBlocker(), func(state tray.State) {
		status.SetTitle(state.Status)
		if state.Error != "" {
			status.SetTitle(fmt.Sprintf("Error: %s", state.Error))
		}

		if state.Active {
			start30s.Disable()
			start1m.Disable()
			start15m.Disable()
			stop.Enable()
			return
		}

		start30s.Enable()
		start1m.Enable()
		start15m.Enable()
		stop.Disable()
	})

	go func() {
		for {
			select {
			case <-start30s.ClickedCh:
				if err := controller.Start(30 * time.Second); err != nil {
					log.Print(err)
				}
			case <-start1m.ClickedCh:
				if err := controller.Start(time.Minute); err != nil {
					log.Print(err)
				}
			case <-start15m.ClickedCh:
				if err := controller.Start(15 * time.Minute); err != nil {
					log.Print(err)
				}
			case <-stop.ClickedCh:
				if err := controller.Stop(); err != nil {
					log.Print(err)
				}
			case <-quit.ClickedCh:
				if err := controller.Stop(); err != nil {
					log.Print(err)
				}
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {}
