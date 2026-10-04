//go:build darwin

#import <Cocoa/Cocoa.h>
#include "graph_darwin.h"

extern void onTrafficGraphClosed(void);

#define HISTORY_CAPACITY 60

@interface TrafficGraphView : NSView
@property (nonatomic, assign) uint64_t currentUpload;
@property (nonatomic, assign) uint64_t currentDownload;
@property (nonatomic, assign) uint64_t peakUpload;
@property (nonatomic, assign) uint64_t peakDownload;
- (void)addTrafficDataWithUpload:(uint64_t)up download:(uint64_t)down;
@end

@implementation TrafficGraphView {
    uint64_t _uploadHistory[HISTORY_CAPACITY];
    uint64_t _downloadHistory[HISTORY_CAPACITY];
    int _count;
}

- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        memset(_uploadHistory, 0, sizeof(_uploadHistory));
        memset(_downloadHistory, 0, sizeof(_downloadHistory));
        _count = 0;
    }
    return self;
}

- (BOOL)isFlipped {
    return YES;
}

- (void)addTrafficDataWithUpload:(uint64_t)up download:(uint64_t)down {
    self.currentUpload = up;
    self.currentDownload = down;

    if (up > self.peakUpload) {
        self.peakUpload = up;
    }
    if (down > self.peakDownload) {
        self.peakDownload = down;
    }

    if (_count < HISTORY_CAPACITY) {
        _uploadHistory[_count] = up;
        _downloadHistory[_count] = down;
        _count++;
    } else {
        memmove(_uploadHistory, _uploadHistory + 1, sizeof(uint64_t) * (HISTORY_CAPACITY - 1));
        memmove(_downloadHistory, _downloadHistory + 1, sizeof(uint64_t) * (HISTORY_CAPACITY - 1));
        _uploadHistory[HISTORY_CAPACITY - 1] = up;
        _downloadHistory[HISTORY_CAPACITY - 1] = down;
    }

    [self setNeedsDisplay:YES];
}

static NSString* formatSpeedObjC(uint64_t bytesPerSec) {
    const uint64_t kb = 1024;
    const uint64_t mb = 1024 * kb;
    const uint64_t gb = 1024 * mb;
    if (bytesPerSec >= gb) {
        return [NSString stringWithFormat:@"%.1f GB/s", (double)bytesPerSec / (double)gb];
    } else if (bytesPerSec >= mb) {
        return [NSString stringWithFormat:@"%.1f MB/s", (double)bytesPerSec / (double)mb];
    } else {
        return [NSString stringWithFormat:@"%llu KB/s", bytesPerSec / kb];
    }
}

- (void)drawRect:(NSRect)dirtyRect {
    [super drawRect:dirtyRect];

    NSRect bounds = self.bounds;

    // Header Legend
    // Upload Dot & Text
    NSRect upDotRect = NSMakeRect(16, 16, 8, 8);
    NSBezierPath *upDot = [NSBezierPath bezierPathWithOvalInRect:upDotRect];
    NSColor *upColor = [NSColor colorWithRed:1.0 green:0.60 blue:0.0 alpha:1.0];
    [upColor setFill];
    [upDot fill];

    NSString *upStr = [NSString stringWithFormat:@"Upload:   ↑ %@  (Peak: %@)",
                       formatSpeedObjC(self.currentUpload), formatSpeedObjC(self.peakUpload)];
    NSDictionary *upAttrs = @{
        NSFontAttributeName: [NSFont systemFontOfSize:11 weight:NSFontWeightMedium],
        NSForegroundColorAttributeName: [NSColor colorWithCalibratedWhite:0.92 alpha:1.0]
    };
    [upStr drawAtPoint:NSMakePoint(30, 13) withAttributes:upAttrs];

    // Download Dot & Text
    NSRect downDotRect = NSMakeRect(16, 33, 8, 8);
    NSBezierPath *downDot = [NSBezierPath bezierPathWithOvalInRect:downDotRect];
    NSColor *downColor = [NSColor colorWithRed:0.0 green:0.75 blue:1.0 alpha:1.0];
    [downColor setFill];
    [downDot fill];

    NSString *downStr = [NSString stringWithFormat:@"Download: ↓ %@  (Peak: %@)",
                         formatSpeedObjC(self.currentDownload), formatSpeedObjC(self.peakDownload)];
    NSDictionary *downAttrs = @{
        NSFontAttributeName: [NSFont systemFontOfSize:11 weight:NSFontWeightMedium],
        NSForegroundColorAttributeName: [NSColor colorWithCalibratedWhite:0.92 alpha:1.0]
    };
    [downStr drawAtPoint:NSMakePoint(30, 30) withAttributes:downAttrs];

    // Plot Dimensions
    CGFloat leftMargin = 58;
    CGFloat rightMargin = 16;
    CGFloat topMargin = 54;
    CGFloat bottomMargin = 26;

    NSRect plotRect = NSMakeRect(leftMargin, topMargin,
                                 bounds.size.width - leftMargin - rightMargin,
                                 bounds.size.height - topMargin - bottomMargin);

    if (plotRect.size.width < 50 || plotRect.size.height < 30) {
        return;
    }

    // Plot Background
    NSBezierPath *plotBg = [NSBezierPath bezierPathWithRoundedRect:plotRect xRadius:4 yRadius:4];
    [[NSColor colorWithCalibratedWhite:0.0 alpha:0.35] setFill];
    [plotBg fill];
    [[NSColor colorWithCalibratedWhite:1.0 alpha:0.12] setStroke];
    [plotBg setLineWidth:1.0];
    [plotBg stroke];

    // Determine scale
    uint64_t maxVal = 0;
    for (int i = 0; i < HISTORY_CAPACITY; i++) {
        if (_uploadHistory[i] > maxVal) maxVal = _uploadHistory[i];
        if (_downloadHistory[i] > maxVal) maxVal = _downloadHistory[i];
    }
    const uint64_t minScale = 100 * 1024; // 100 KB/s floor
    if (maxVal < minScale) {
        maxVal = minScale;
    }
    double scaleMax = (double)maxVal * 1.15; // 15% headroom

    // Grid lines & Y-axis labels
    NSDictionary *axisAttrs = @{
        NSFontAttributeName: [NSFont monospacedDigitSystemFontOfSize:9 weight:NSFontWeightRegular],
        NSForegroundColorAttributeName: [NSColor colorWithCalibratedWhite:0.55 alpha:1.0]
    };

    CGFloat plotBottom = plotRect.origin.y + plotRect.size.height;
    CGFloat plotTop = plotRect.origin.y;

    // 100% line & label
    NSString *scaleTopStr = formatSpeedObjC((uint64_t)scaleMax);
    NSSize topStrSize = [scaleTopStr sizeWithAttributes:axisAttrs];
    [scaleTopStr drawAtPoint:NSMakePoint(leftMargin - topStrSize.width - 6, plotTop - 2) withAttributes:axisAttrs];

    // 50% line & label
    CGFloat midY = plotTop + plotRect.size.height * 0.5;
    NSString *scaleMidStr = formatSpeedObjC((uint64_t)(scaleMax * 0.5));
    NSSize midStrSize = [scaleMidStr sizeWithAttributes:axisAttrs];
    [scaleMidStr drawAtPoint:NSMakePoint(leftMargin - midStrSize.width - 6, midY - 6) withAttributes:axisAttrs];

    NSBezierPath *midGrid = [NSBezierPath bezierPath];
    [midGrid moveToPoint:NSMakePoint(plotRect.origin.x, midY)];
    [midGrid lineToPoint:NSMakePoint(plotRect.origin.x + plotRect.size.width, midY)];
    CGFloat dash[2] = {3.0, 3.0};
    [midGrid setLineDash:dash count:2 phase:0.0];
    [[NSColor colorWithCalibratedWhite:1.0 alpha:0.10] setStroke];
    [midGrid stroke];

    // 0 line & label
    NSString *scaleBotStr = @"0 KB/s";
    NSSize botStrSize = [scaleBotStr sizeWithAttributes:axisAttrs];
    [scaleBotStr drawAtPoint:NSMakePoint(leftMargin - botStrSize.width - 6, plotBottom - 9) withAttributes:axisAttrs];

    // X-axis time markings
    NSDictionary *timeAttrs = @{
        NSFontAttributeName: [NSFont systemFontOfSize:9 weight:NSFontWeightRegular],
        NSForegroundColorAttributeName: [NSColor colorWithCalibratedWhite:0.45 alpha:1.0]
    };
    [@"-60s" drawAtPoint:NSMakePoint(plotRect.origin.x, plotBottom + 5) withAttributes:timeAttrs];
    [@"-30s" drawAtPoint:NSMakePoint(plotRect.origin.x + plotRect.size.width * 0.5 - 10, plotBottom + 5) withAttributes:timeAttrs];
    [@"Now" drawAtPoint:NSMakePoint(plotRect.origin.x + plotRect.size.width - 22, plotBottom + 5) withAttributes:timeAttrs];

    // Compute curve points
    NSPoint downPoints[HISTORY_CAPACITY];
    NSPoint upPoints[HISTORY_CAPACITY];

    for (int i = 0; i < HISTORY_CAPACITY; i++) {
        CGFloat x = plotRect.origin.x + ((CGFloat)i / (CGFloat)(HISTORY_CAPACITY - 1)) * plotRect.size.width;

        double downRatio = (double)_downloadHistory[i] / scaleMax;
        if (downRatio > 1.0) downRatio = 1.0;
        CGFloat yDown = plotBottom - downRatio * plotRect.size.height;
        downPoints[i] = NSMakePoint(x, yDown);

        double upRatio = (double)_uploadHistory[i] / scaleMax;
        if (upRatio > 1.0) upRatio = 1.0;
        CGFloat yUp = plotBottom - upRatio * plotRect.size.height;
        upPoints[i] = NSMakePoint(x, yUp);
    }

    // 1. Draw Download Series (Area + Line)
    NSBezierPath *downArea = [NSBezierPath bezierPath];
    [downArea moveToPoint:NSMakePoint(downPoints[0].x, plotBottom)];
    for (int i = 0; i < HISTORY_CAPACITY; i++) {
        [downArea lineToPoint:downPoints[i]];
    }
    [downArea lineToPoint:NSMakePoint(downPoints[HISTORY_CAPACITY - 1].x, plotBottom)];
    [downArea closePath];

    [NSGraphicsContext saveGraphicsState];
    [plotBg addClip]; // Clip inside rounded plot rect
    [downArea addClip];
    NSColor *downGradTop = [NSColor colorWithRed:0.0 green:0.75 blue:1.0 alpha:0.30];
    NSColor *downGradBot = [NSColor colorWithRed:0.0 green:0.75 blue:1.0 alpha:0.02];
    NSGradient *downGrad = [[NSGradient alloc] initWithStartingColor:downGradBot endingColor:downGradTop];
    [downGrad drawInRect:plotRect angle:90.0];
    [NSGraphicsContext restoreGraphicsState];

    NSBezierPath *downLine = [NSBezierPath bezierPath];
    [downLine setLineWidth:1.8];
    [downLine setLineJoinStyle:NSLineJoinStyleRound];
    [downLine setLineCapStyle:NSLineCapStyleRound];
    [downLine moveToPoint:downPoints[0]];
    for (int i = 1; i < HISTORY_CAPACITY; i++) {
        [downLine lineToPoint:downPoints[i]];
    }
    [NSGraphicsContext saveGraphicsState];
    [plotBg addClip];
    [downColor setStroke];
    [downLine stroke];

    // Download pulse dot
    NSPoint lastDown = downPoints[HISTORY_CAPACITY - 1];
    NSBezierPath *dDot = [NSBezierPath bezierPathWithOvalInRect:NSMakeRect(lastDown.x - 2.5, lastDown.y - 2.5, 5.0, 5.0)];
    [downColor setFill];
    [dDot fill];
    [NSGraphicsContext restoreGraphicsState];

    // 2. Draw Upload Series (Area + Line)
    NSBezierPath *upArea = [NSBezierPath bezierPath];
    [upArea moveToPoint:NSMakePoint(upPoints[0].x, plotBottom)];
    for (int i = 0; i < HISTORY_CAPACITY; i++) {
        [upArea lineToPoint:upPoints[i]];
    }
    [upArea lineToPoint:NSMakePoint(upPoints[HISTORY_CAPACITY - 1].x, plotBottom)];
    [upArea closePath];

    [NSGraphicsContext saveGraphicsState];
    [plotBg addClip];
    [upArea addClip];
    NSColor *upGradTop = [NSColor colorWithRed:1.0 green:0.60 blue:0.0 alpha:0.25];
    NSColor *upGradBot = [NSColor colorWithRed:1.0 green:0.60 blue:0.0 alpha:0.02];
    NSGradient *upGrad = [[NSGradient alloc] initWithStartingColor:upGradBot endingColor:upGradTop];
    [upGrad drawInRect:plotRect angle:90.0];
    [NSGraphicsContext restoreGraphicsState];

    NSBezierPath *upLine = [NSBezierPath bezierPath];
    [upLine setLineWidth:1.8];
    [upLine setLineJoinStyle:NSLineJoinStyleRound];
    [upLine setLineCapStyle:NSLineCapStyleRound];
    [upLine moveToPoint:upPoints[0]];
    for (int i = 1; i < HISTORY_CAPACITY; i++) {
        [upLine lineToPoint:upPoints[i]];
    }
    [NSGraphicsContext saveGraphicsState];
    [plotBg addClip];
    [upColor setStroke];
    [upLine stroke];

    // Upload pulse dot
    NSPoint lastUp = upPoints[HISTORY_CAPACITY - 1];
    NSBezierPath *uDot = [NSBezierPath bezierPathWithOvalInRect:NSMakeRect(lastUp.x - 2.5, lastUp.y - 2.5, 5.0, 5.0)];
    [upColor setFill];
    [uDot fill];
    [NSGraphicsContext restoreGraphicsState];
}

@end

@interface TrafficGraphWindowDelegate : NSObject <NSWindowDelegate>
@end

@implementation TrafficGraphWindowDelegate
- (void)windowWillClose:(NSNotification *)notification {
    onTrafficGraphClosed();
}
@end

static NSPanel *sharedGraphPanel = nil;
static TrafficGraphView *sharedGraphView = nil;
static TrafficGraphWindowDelegate *sharedGraphDelegate = nil;

static void createTrafficGraphPanel(void) {
    if (sharedGraphPanel) return;

    NSRect frame = NSMakeRect(200, 200, 440, 260);
    sharedGraphPanel = [[NSPanel alloc] initWithContentRect:frame
                                                  styleMask:NSWindowStyleMaskTitled |
                                                            NSWindowStyleMaskClosable |
                                                            NSWindowStyleMaskResizable |
                                                            NSWindowStyleMaskUtilityWindow
                                                    backing:NSBackingStoreBuffered
                                                      defer:NO];
    [sharedGraphPanel setTitle:@"Network Traffic Graph"];
    [sharedGraphPanel setLevel:NSFloatingWindowLevel];
    [sharedGraphPanel setReleasedWhenClosed:NO];
    [sharedGraphPanel setMinSize:NSMakeSize(340, 200)];
    [sharedGraphPanel setMovableByWindowBackground:YES];

    sharedGraphDelegate = [[TrafficGraphWindowDelegate alloc] init];
    [sharedGraphPanel setDelegate:sharedGraphDelegate];

    NSVisualEffectView *vibrantView = [[NSVisualEffectView alloc] initWithFrame:[sharedGraphPanel.contentView bounds]];
    vibrantView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    vibrantView.material = NSVisualEffectMaterialHUDWindow;
    vibrantView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
    vibrantView.state = NSVisualEffectStateActive;
    [sharedGraphPanel setContentView:vibrantView];

    sharedGraphView = [[TrafficGraphView alloc] initWithFrame:vibrantView.bounds];
    sharedGraphView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    [vibrantView addSubview:sharedGraphView];
}

void toggleTrafficGraphWindow(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        createTrafficGraphPanel();
        if ([sharedGraphPanel isVisible]) {
            [sharedGraphPanel orderOut:nil];
            onTrafficGraphClosed();
        } else {
            [sharedGraphPanel makeKeyAndOrderFront:nil];
            [NSApp activateIgnoringOtherApps:YES];
        }
    });
}

void updateTrafficGraph(uint64_t upSpeed, uint64_t downSpeed) {
    dispatch_async(dispatch_get_main_queue(), ^{
        createTrafficGraphPanel();
        [sharedGraphView addTrafficDataWithUpload:upSpeed download:downSpeed];
    });
}

int isTrafficGraphVisible(void) {
    if (!sharedGraphPanel) return 0;
    return [sharedGraphPanel isVisible] ? 1 : 0;
}
