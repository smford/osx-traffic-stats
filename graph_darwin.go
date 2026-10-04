//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include <stdlib.h>
#include "graph_darwin.h"
*/
import "C"
import (
	"sync/atomic"
	"unsafe"
)

var isGraphVisible atomic.Bool

//export onTrafficGraphClosed
func onTrafficGraphClosed() {
	isGraphVisible.Store(false)
	if mToggleGraph != nil {
		mToggleGraph.SetTitle("Show Traffic Graph")
	}
}

func toggleTrafficGraph() {
	C.toggleTrafficGraphWindow()
	if isTrafficGraphVisible() {
		isGraphVisible.Store(false)
		if mToggleGraph != nil {
			mToggleGraph.SetTitle("Show Traffic Graph")
		}
	} else {
		isGraphVisible.Store(true)
		if mToggleGraph != nil {
			mToggleGraph.SetTitle("Hide Traffic Graph")
		}
	}
}

func updateGraph(upSpeed, downSpeed uint64) {
	summary := FormatTopAppsSummary(3)
	cSummary := C.CString(summary)
	defer C.free(unsafe.Pointer(cSummary))
	C.updateTrafficGraph(C.uint64_t(upSpeed), C.uint64_t(downSpeed), cSummary)
}

func isTrafficGraphVisible() bool {
	return C.isTrafficGraphVisible() != 0
}

func resetPeaks() {
	C.resetTrafficGraphPeaks()
}

func setGraphOpacity(opacity float64) {
	C.setTrafficGraphOpacity(C.double(opacity))
}

func setGraphAlwaysOnTop(alwaysOnTop bool) {
	var val C.int
	if alwaysOnTop {
		val = 1
	}
	C.setTrafficGraphAlwaysOnTop(val)
}

func snapGraphWindow(corner int) {
	C.snapTrafficGraphWindow(C.int(corner))
}
