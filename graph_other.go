//go:build !darwin

package main

func toggleTrafficGraph() {}

func updateGraph(upSpeed, downSpeed uint64) {}

func isTrafficGraphVisible() bool {
	return false
}
