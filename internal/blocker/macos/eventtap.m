#import "eventtap.h"
#import <AppKit/AppKit.h>
#import <ApplicationServices/ApplicationServices.h>
#import <CoreFoundation/Corefoundation.h>
#import <IOKit/hidsystem/ev_keymap.h>

static CFMachPortRef eventTap = NULL;
static CFRunLoopSourceRef runLoopSource = NULL;
static CFRunLoopRef eventLoop = NULL;
static dispatch_queue_t eventQueue = NULL;

static const CGEventType systemDefinedEventType = (CGEventType)NX_SYSDEFINED;

static CGEventRef keyboardCallBack(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
	switch(type) {
		case kCGEventKeyDown:
		case kCGEventKeyUp:
		case kCGEventFlagsChanged:
			return NULL;
		case NX_SYSDEFINED: {
			NSEvent *nsEvent = [NSEvent eventWithCGEvent:event];
			if (nsEvent != nil && [nsEvent subtype] == NX_SUBTYPE_AUX_CONTROL_BUTTONS) {
				return NULL;
			}
			return event;
		}
		default:
			return event;
	}
}

int StartKeyboardBlock(void) {
	if (eventTap != NULL) {
		return 1;
	}

	dispatch_semaphore_t started = dispatch_semaphore_create(0);
	__block int ok = 0;

	if (eventQueue == NULL) {
		eventQueue = dispatch_queue_create("go-keyboard-cleaner.eventtap", DISPATCH_QUEUE_SERIAL);
	}

	dispatch_async(eventQueue, ^{
		CGEventMask mask = CGEventMaskBit(kCGEventKeyDown) | CGEventMaskBit(kCGEventKeyUp) | CGEventMaskBit(kCGEventFlagsChanged) | CGEventMaskBit(systemDefinedEventType);

		eventTap = CGEventTapCreate(kCGSessionEventTap, kCGHeadInsertEventTap, kCGEventTapOptionDefault, mask, keyboardCallBack, NULL);
		if (eventTap == NULL) {
			dispatch_semaphore_signal(started);
			return;
		}

		runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
		eventLoop = CFRunLoopGetCurrent();

		CFRunLoopAddSource(eventLoop, runLoopSource, kCFRunLoopCommonModes);
		CGEventTapEnable(eventTap, true);

		ok = 1;
		dispatch_semaphore_signal(started);

		CFRunLoopRun();
	});

	dispatch_semaphore_wait(started, DISPATCH_TIME_FOREVER);
	return ok;
}

void StopKeyboardBlock(void) {
	if (eventTap == NULL) {
		return;
	}

	CGEventTapEnable(eventTap, false);

	if (eventLoop != NULL && runLoopSource != NULL) {
		CFRunLoopRemoveSource(eventLoop, runLoopSource, kCFRunLoopCommonModes);
	}

	if (eventLoop != NULL) {
		CFRunLoopStop(eventLoop);
	}

	if (runLoopSource != NULL) {
		CFRelease(runLoopSource);
		runLoopSource = NULL;
	}

	if (eventTap != NULL) {
		CFRelease(eventTap);
		eventTap = NULL;
	}

	eventLoop = NULL;
}

int AccessibilityTrusted(int prompt) {
	const void *keys[] = { kAXTrustedCheckOptionPrompt };
	const void *values[] = { prompt ? kCFBooleanTrue : kCFBooleanFalse };
	CFDictionaryRef options = CFDictionaryCreate(
		kCFAllocatorDefault,
		keys,
		values,
		1,
		&kCFCopyStringDictionaryKeyCallBacks,
		&kCFTypeDictionaryValueCallBacks
	);

	Boolean trusted = AXIsProcessTrustedWithOptions(options);
	CFRelease(options);

	return trusted ? 1 : 0;
}
