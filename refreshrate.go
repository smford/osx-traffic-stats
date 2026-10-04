package main

import (
	"sync/atomic"
	"time"
)

type RefreshRate int

const (
	RefreshFast    RefreshRate = 500  // 0.5s
	RefreshNormal  RefreshRate = 1000 // 1.0s (default)
	RefreshBattery RefreshRate = 2000 // 2.0s
	RefreshEco     RefreshRate = 5000 // 5.0s
)

var currentRefreshRate atomic.Int64

func init() {
	currentRefreshRate.Store(int64(RefreshNormal))
}

func getRefreshRate() RefreshRate {
	val := currentRefreshRate.Load()
	if val <= 0 {
		return RefreshNormal
	}
	return RefreshRate(val)
}

func setRefreshRate(r RefreshRate) {
	switch r {
	case RefreshFast, RefreshNormal, RefreshBattery, RefreshEco:
		currentRefreshRate.Store(int64(r))
	default:
		currentRefreshRate.Store(int64(RefreshNormal))
	}
}

func getRefreshDuration() time.Duration {
	return time.Duration(getRefreshRate()) * time.Millisecond
}

// CalculateSpeedPerSecond normalizes delta bytes over the elapsed duration to bytes per second.
func CalculateSpeedPerSecond(deltaBytes uint64, elapsedSec float64) uint64 {
	if elapsedSec <= 0 {
		return deltaBytes
	}
	return uint64(float64(deltaBytes) / elapsedSec)
}
