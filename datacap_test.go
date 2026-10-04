package main

import (
	"testing"
)

func TestDataCapMonitoring(t *testing.T) {
	var notifiedCount int
	notifyDataCap = func(currentBytes, capBytes uint64) {
		notifiedCount++
	}
	defer func() {
		notifyDataCap = sendDataCapNotification
	}()

	// Disabled
	setDataCap(0)
	if checkDataCap(100 * 1024 * 1024 * 1024) {
		t.Error("checkDataCap should return false when disabled")
	}
	if notifiedCount != 0 {
		t.Errorf("notifiedCount = %d, want 0 when disabled", notifiedCount)
	}

	// 1 GB threshold
	const oneGB = 1024 * 1024 * 1024
	setDataCap(oneGB)

	// Below cap
	if checkDataCap(500 * 1024 * 1024) {
		t.Error("checkDataCap should return false when below cap")
	}
	if notifiedCount != 0 {
		t.Errorf("notifiedCount = %d, want 0 when below cap", notifiedCount)
	}

	// Reached cap
	if !checkDataCap(oneGB) {
		t.Error("checkDataCap should return true when reaching cap")
	}
	if notifiedCount != 1 {
		t.Errorf("notifiedCount = %d, want 1 when reaching cap", notifiedCount)
	}

	// Subsequent checks should return false (already alerted)
	if checkDataCap(oneGB + 1024) {
		t.Error("checkDataCap should return false once already alerted")
	}
	if notifiedCount != 1 {
		t.Errorf("notifiedCount = %d, want 1 after already alerted", notifiedCount)
	}

	// Reset alert
	resetDataCapAlert()
	if !checkDataCap(oneGB + 2048) {
		t.Error("checkDataCap should return true after alert reset")
	}
	if notifiedCount != 2 {
		t.Errorf("notifiedCount = %d, want 2 after alert reset", notifiedCount)
	}

	// Cleanup
	setDataCap(0)
}
