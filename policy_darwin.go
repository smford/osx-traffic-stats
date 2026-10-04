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

#define MINI_SAMPLES 16
static uint64_t miniUpHistory[MINI_SAMPLES] = {0};
static uint64_t miniDownHistory[MINI_SAMPLES] = {0};
static int miniHistoryIndex = 0;

void updateMenuBarGraph(uint64_t upSpeed, uint64_t downSpeed, int iconMode) {
	dispatch_async(dispatch_get_main_queue(), ^{
		id delegate = [NSApp delegate];
		if (!delegate) return;
		Ivar ivar = class_getInstanceVariable([delegate class], "statusItem");
		if (!ivar) return;
		NSStatusItem *statusItem = object_getIvar(delegate, ivar);
		if (!statusItem || !statusItem.button) return;

		if (iconMode == 0) {
			if (statusItem.button.image != nil) {
				[statusItem.button setImage:nil];
				[statusItem.button setImagePosition:NSNoImage];
				[statusItem setLength:172.0];
			}
			return;
		}

		miniUpHistory[miniHistoryIndex] = upSpeed;
		miniDownHistory[miniHistoryIndex] = downSpeed;
		miniHistoryIndex = (miniHistoryIndex + 1) % MINI_SAMPLES;

		uint64_t maxSpeed = 10 * 1024;
		for (int i = 0; i < MINI_SAMPLES; i++) {
			if (miniUpHistory[i] > maxSpeed) maxSpeed = miniUpHistory[i];
			if (miniDownHistory[i] > maxSpeed) maxSpeed = miniDownHistory[i];
		}

		CGFloat width = 28.0;
		CGFloat height = 18.0;

		NSImage *img = [NSImage imageWithSize:NSMakeSize(width, height) flipped:NO drawingHandler:^BOOL(NSRect dstRect) {
			NSBezierPath *bg = [NSBezierPath bezierPathWithRoundedRect:dstRect xRadius:2.0 yRadius:2.0];
			[[NSColor colorWithCalibratedWhite:0.2 alpha:0.35] setFill];
			[bg fill];

			CGFloat halfH = (height - 3.0) / 2.0;
			CGFloat step = (width - 2.0) / (CGFloat)(MINI_SAMPLES - 1);

			// Upload path (Amber/Orange)
			NSBezierPath *upPath = [NSBezierPath bezierPath];
			[[NSColor colorWithCalibratedRed:1.0 green:0.65 blue:0.1 alpha:0.95] setStroke];
			[upPath setLineWidth:1.2];

			for (int i = 0; i < MINI_SAMPLES; i++) {
				int idx = (miniHistoryIndex + i) % MINI_SAMPLES;
				CGFloat x = 1.0 + i * step;
				CGFloat ratio = (CGFloat)miniUpHistory[idx] / (CGFloat)maxSpeed;
				if (ratio > 1.0) ratio = 1.0;
				CGFloat y = (height / 2.0) + (ratio * halfH);
				if (i == 0) [upPath moveToPoint:NSMakePoint(x, y)];
				else [upPath lineToPoint:NSMakePoint(x, y)];
			}
			[upPath stroke];

			// Download path (Cyan/Blue)
			NSBezierPath *downPath = [NSBezierPath bezierPath];
			[[NSColor colorWithCalibratedRed:0.2 green:0.8 blue:1.0 alpha:0.95] setStroke];
			[downPath setLineWidth:1.2];

			for (int i = 0; i < MINI_SAMPLES; i++) {
				int idx = (miniHistoryIndex + i) % MINI_SAMPLES;
				CGFloat x = 1.0 + i * step;
				CGFloat ratio = (CGFloat)miniDownHistory[idx] / (CGFloat)maxSpeed;
				if (ratio > 1.0) ratio = 1.0;
				CGFloat y = (height / 2.0) - (ratio * halfH);
				if (i == 0) [downPath moveToPoint:NSMakePoint(x, y)];
				else [downPath lineToPoint:NSMakePoint(x, y)];
			}
			[downPath stroke];

			return YES;
		}];
		[img setTemplate:NO];

		[statusItem.button setImage:img];
		if (iconMode == 2) {
			[statusItem.button setImagePosition:NSImageOnly];
			[statusItem setLength:34.0];
		} else {
			[statusItem.button setImagePosition:NSImageLeading];
			[statusItem setLength:206.0];
		}
	});
}

static BOOL aboutBoxOpen = NO;

void showAboutBox(const char *title, const char *message, const char *url, const void *iconData, int iconLen) {
	NSString *nsTitle = [NSString stringWithUTF8String:title];
	NSString *nsMessage = [NSString stringWithUTF8String:message];
	NSString *nsUrl = [NSString stringWithUTF8String:url];
	NSData *nsIconData = (iconData != NULL && iconLen > 0) ? [NSData dataWithBytes:iconData length:iconLen] : nil;

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

			if (nsIconData) {
				NSImage *iconImg = [[NSImage alloc] initWithData:nsIconData];
				if (iconImg) {
					[iconImg setSize:NSMakeSize(64.0, 64.0)];
					[alert setIcon:iconImg];
				}
			}

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

	var iconPtr unsafe.Pointer
	var iconLen C.int
	if len(appLogoSVG) > 0 {
		iconPtr = unsafe.Pointer(&appLogoSVG[0])
		iconLen = C.int(len(appLogoSVG))
	}

	C.showAboutBox(cTitle, cMessage, cURL, iconPtr, iconLen)
}

func updateMenuBarGraph(upSpeed, downSpeed uint64, mode MenuBarIconMode) {
	C.updateMenuBarGraph(C.uint64_t(upSpeed), C.uint64_t(downSpeed), C.int(mode))
}
