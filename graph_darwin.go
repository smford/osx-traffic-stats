//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa

#include "graph_darwin.h"
*/
import "C"
import "sync/atomic"

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
	C.updateTrafficGraph(C.uint64_t(upSpeed), C.uint64_t(downSpeed))
}

func isTrafficGraphVisible() bool {
	return C.isTrafficGraphVisible() != 0
}

func resetPeaks() {
	C.resetTrafficGraphPeaks()
}
