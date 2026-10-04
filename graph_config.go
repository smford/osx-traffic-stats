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
