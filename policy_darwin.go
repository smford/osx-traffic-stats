//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>
#include <stdlib.h>
#include <math.h>

void setAccessoryPolicy(void) {
	[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
}

void configureFixedStatusItem(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		id delegate = [NSApp delegate];
		if (!delegate) return;

		Ivar ivar = class_getInstanceVariable([delegate class], "statusItem");
		if (!ivar) return;

		NSStatusItem *statusItem = object_getIvar(delegate, ivar);
		if (!statusItem) return;

		NSFont *font = [NSFont monospacedDigitSystemFontOfSize:12.0 weight:NSFontWeightRegular];
		NSString *sample = @"↑  999 KB/s  ↓  999 KB/s";
		NSDictionary *attrs = @{ NSFontAttributeName: font };
		NSSize size = [sample sizeWithAttributes:attrs];

		CGFloat itemLength = ceil(size.width + 16.0);
		[statusItem setLength:itemLength];

		if (statusItem.button) {
			statusItem.button.font = font;
			statusItem.button.alignment = NSTextAlignmentCenter;
		}
	});
}

static BOOL aboutBoxOpen = NO;

void showAboutBox(const char *title, const char *message, const char *url) {
	NSString *nsTitle = [NSString stringWithUTF8String:title];
	NSString *nsMessage = [NSString stringWithUTF8String:message];
	NSString *nsUrl = [NSString stringWithUTF8String:url];

	dispatch_async(dispatch_get_main_queue(), ^{
		if (aboutBoxOpen) {
			[NSApp activateIgnoringOtherApps:YES];
			return;
		}
		aboutBoxOpen = YES;
		@try {
			NSAlert *alert = [[NSAlert alloc] init];
			[alert setMessageText:nsTitle];
			[alert setInformativeText:nsMessage];
			[alert setAlertStyle:NSAlertStyleInformational];
			[alert addButtonWithTitle:@"OK"];
			[alert addButtonWithTitle:@"View on GitHub"];

			[NSApp activateIgnoringOtherApps:YES];
			[[alert window] setLevel:NSFloatingWindowLevel];
			NSModalResponse response = [alert runModal];
			if (response == NSAlertSecondButtonReturn) {
				[[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:nsUrl]];
			}
		} @finally {
			aboutBoxOpen = NO;
		}
	});
}
*/
import "C"
import "unsafe"

func setAccessoryPolicy() {
	C.setAccessoryPolicy()
}

func configureFixedStatusItem() {
	C.configureFixedStatusItem()
}

func showAboutBox() {
	cTitle := C.CString(appName)
	defer C.free(unsafe.Pointer(cTitle))

	cMessage := C.CString(aboutMessage())
	defer C.free(unsafe.Pointer(cMessage))

	cURL := C.CString(githubURL)
	defer C.free(unsafe.Pointer(cURL))

	C.showAboutBox(cTitle, cMessage, cURL)
}
