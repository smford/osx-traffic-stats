//go:build !darwin

package main

func toggleTrafficGraph() {}

func updateGraph(upSpeed, downSpeed uint64) {}

func isTrafficGraphVisible() bool {
	return false
}

func resetPeaks() {}

func setGraphOpacity(opacity float64) {}

func setGraphAlwaysOnTop(alwaysOnTop bool) {}

func snapGraphWindow(corner int) {}

