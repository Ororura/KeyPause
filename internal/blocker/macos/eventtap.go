package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation -framework AppKit

#include "eventtap.h"
*/
import "C"

func startEventTap() bool {
	return C.StartKeyboardBlock() == 1
}

func stopEventTap() {
	C.StopKeyboardBlock()
}

func accessibilityTrusted(prompt bool) bool {
	var shouldPrompt C.int
	if prompt {
		shouldPrompt = 1
	}

	return C.AccessibilityTrusted(shouldPrompt) == 1
}
