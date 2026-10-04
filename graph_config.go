package main

import (
	"math"
	"sync/atomic"
)

var (
	graphOpacityBits   atomic.Uint64
	graphAlwaysOnTopVal atomic.Bool
)

func init() {
	graphOpacityBits.Store(math.Float64bits(1.0))
	graphAlwaysOnTopVal.Store(true)
}

func getGraphOpacity() float64 {
	bits := graphOpacityBits.Load()
	val := math.Float64frombits(bits)
	if val <= 0.1 || val > 1.0 {
		return 1.0
	}
	return val
}

func setGraphOpacityConfig(val float64) {
	if val < 0.2 {
		val = 0.2
	} else if val > 1.0 {
		val = 1.0
	}
	graphOpacityBits.Store(math.Float64bits(val))
	setGraphOpacity(val)
}

func getGraphAlwaysOnTop() bool {
	return graphAlwaysOnTopVal.Load()
}

func setGraphAlwaysOnTopConfig(val bool) {
	graphAlwaysOnTopVal.Store(val)
	setGraphAlwaysOnTop(val)
}

// ResampleTrafficHistory resamples raw second-by-second history samples into displayPoints buckets
// for a specified timeframe window (in seconds).
func ResampleTrafficHistory(history []uint64, windowSec, displayPoints int) []uint64 {
	if displayPoints <= 0 {
		return nil
	}
	if windowSec <= 0 {
		windowSec = 60
	}
	bucketSize := windowSec / displayPoints
	if bucketSize < 1 {
		bucketSize = 1
	}

	result := make([]uint64, displayPoints)
	count := len(history)
	startIndex := count - windowSec

	for p := 0; p < displayPoints; p++ {
		var maxInBucket uint64
		bStart := startIndex + p*bucketSize
		for b := 0; b < bucketSize; b++ {
			idx := bStart + b
			if idx >= 0 && idx < count {
				val := history[idx]
				if val > maxInBucket {
					maxInBucket = val
				}
			}
		}
		result[p] = maxInBucket
	}
	return result
}

