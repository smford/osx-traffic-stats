package main

import (
	"testing"
)

func TestGraphConfig(t *testing.T) {
	// Defaults
	if getGraphOpacity() != 1.0 {
		t.Errorf("expected default opacity 1.0, got %f", getGraphOpacity())
	}
	if !getGraphAlwaysOnTop() {
		t.Errorf("expected default alwaysOnTop true")
	}

	// Set opacity
	setGraphOpacityConfig(0.85)
	if getGraphOpacity() != 0.85 {
		t.Errorf("expected opacity 0.85, got %f", getGraphOpacity())
	}

	// Clamping
	setGraphOpacityConfig(1.5)
	if getGraphOpacity() != 1.0 {
		t.Errorf("expected clamped opacity 1.0, got %f", getGraphOpacity())
	}

	setGraphOpacityConfig(0.05)
	if getGraphOpacity() != 0.2 {
		t.Errorf("expected clamped opacity 0.2, got %f", getGraphOpacity())
	}

	// Reset opacity
	setGraphOpacityConfig(1.0)

	// Always on top toggle
	setGraphAlwaysOnTopConfig(false)
	if getGraphAlwaysOnTop() {
		t.Errorf("expected alwaysOnTop false")
	}

	setGraphAlwaysOnTopConfig(true)
	if !getGraphAlwaysOnTop() {
		t.Errorf("expected alwaysOnTop true")
	}
}

func TestResampleTrafficHistory(t *testing.T) {
	// 1. Empty history
	empty := ResampleTrafficHistory([]uint64{}, 60, 60)
	for i, v := range empty {
		if v != 0 {
			t.Errorf("expected 0 for empty history at index %d, got %d", i, v)
		}
	}

	// 2. 1m interval with 10 samples (fewer than 60)
	// The 10 samples should be at the rightmost 10 display points (indices 50..59)
	tenSamples := make([]uint64, 10)
	for i := 0; i < 10; i++ {
		tenSamples[i] = uint64(i + 1) // 1, 2, ..., 10
	}
	res1mShort := ResampleTrafficHistory(tenSamples, 60, 60)
	if len(res1mShort) != 60 {
		t.Fatalf("expected 60 points, got %d", len(res1mShort))
	}
	// Indices 0..49 must be 0
	for i := 0; i < 50; i++ {
		if res1mShort[i] != 0 {
			t.Errorf("expected 0 at index %d, got %d", i, res1mShort[i])
		}
	}
	// Indices 50..59 must match the 10 samples (1..10)
	for i := 0; i < 10; i++ {
		if res1mShort[50+i] != uint64(i+1) {
			t.Errorf("expected %d at index %d, got %d", i+1, 50+i, res1mShort[50+i])
		}
	}

	// 3. 1m interval with 100 samples (> 60)
	// Indices 0..59 should correspond to samples 40..99
	hundredSamples := make([]uint64, 100)
	for i := 0; i < 100; i++ {
		hundredSamples[i] = uint64(i + 100)
	}
	res1mLong := ResampleTrafficHistory(hundredSamples, 60, 60)
	for i := 0; i < 60; i++ {
		expected := hundredSamples[40+i]
		if res1mLong[i] != expected {
			t.Errorf("expected %d at index %d, got %d", expected, i, res1mLong[i])
		}
	}

	// 4. 5m interval (300s window into 60 display points, bucketSize = 5)
	// 300 samples with distinct values
	threeHundredSamples := make([]uint64, 300)
	for i := 0; i < 300; i++ {
		threeHundredSamples[i] = uint64(i + 1)
	}
	res5m := ResampleTrafficHistory(threeHundredSamples, 300, 60)
	if len(res5m) != 60 {
		t.Fatalf("expected 60 points, got %d", len(res5m))
	}
	// For each bucket p, max should be the last item in that 5-sample bucket: (p+1)*5
	for p := 0; p < 60; p++ {
		expected := uint64((p + 1) * 5)
		if res5m[p] != expected {
			t.Errorf("bucket %d: expected %d, got %d", p, expected, res5m[p])
		}
	}

	// 5. 15m interval (900s window into 60 display points, bucketSize = 15)
	nineHundredSamples := make([]uint64, 900)
	for i := 0; i < 900; i++ {
		nineHundredSamples[i] = uint64(i + 1)
	}
	res15m := ResampleTrafficHistory(nineHundredSamples, 900, 60)
	if len(res15m) != 60 {
		t.Fatalf("expected 60 points, got %d", len(res15m))
	}
	for p := 0; p < 60; p++ {
		expected := uint64((p + 1) * 15)
		if res15m[p] != expected {
			t.Errorf("15m bucket %d: expected %d, got %d", p, expected, res15m[p])
		}
	}
}

