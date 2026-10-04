package main

import (
	"testing"
	"time"
)

func TestRefreshRateConfig(t *testing.T) {
	setRefreshRate(RefreshFast)
	if getRefreshRate() != RefreshFast {
		t.Errorf("expected RefreshFast, got %d", getRefreshRate())
	}
	if getRefreshDuration() != 500*time.Millisecond {
		t.Errorf("expected 500ms duration, got %v", getRefreshDuration())
	}

	setRefreshRate(RefreshBattery)
	if getRefreshRate() != RefreshBattery {
		t.Errorf("expected RefreshBattery, got %d", getRefreshRate())
	}
	if getRefreshDuration() != 2000*time.Millisecond {
		t.Errorf("expected 2000ms duration, got %v", getRefreshDuration())
	}

	// Invalid value falls back to normal
	setRefreshRate(RefreshRate(99999))
	if getRefreshRate() != RefreshNormal {
		t.Errorf("expected fallback to RefreshNormal, got %d", getRefreshRate())
	}

	setRefreshRate(RefreshNormal)
}

func TestCalculateSpeedPerSecond(t *testing.T) {
	// 1 MB over 1.0s -> 1 MB/s
	speed := CalculateSpeedPerSecond(1048576, 1.0)
	if speed != 1048576 {
		t.Errorf("expected 1048576, got %d", speed)
	}

	// 500 KB over 0.5s -> 1,000,000 B/s
	speed = CalculateSpeedPerSecond(500000, 0.5)
	if speed != 1000000 {
		t.Errorf("expected 1000000, got %d", speed)
	}

	// 10 MB over 2.0s -> 5 MB/s
	speed = CalculateSpeedPerSecond(10000000, 2.0)
	if speed != 5000000 {
		t.Errorf("expected 5000000, got %d", speed)
	}

	// 0 elapsed fallback
	speed = CalculateSpeedPerSecond(5000, 0.0)
	if speed != 5000 {
		t.Errorf("expected 5000, got %d", speed)
	}
}
