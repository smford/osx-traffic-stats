package main

import (
	"testing"
)

func TestDataCapMonitoring(t *testing.T) {
	// Disabled
	setDataCap(0)
	if checkDataCap(100 * 1024 * 1024 * 1024) {
		t.Error("checkDataCap should return false when disabled")
	}

	// 1 GB threshold
	const oneGB = 1024 * 1024 * 1024
	setDataCap(oneGB)

	// Below cap
	if checkDataCap(500 * 1024 * 1024) {
		t.Error("checkDataCap should return false when below cap")
	}

	// Reached cap
	if !checkDataCap(oneGB) {
		t.Error("checkDataCap should return true when reaching cap")
	}

	// Subsequent checks should return false (already alerted)
	if checkDataCap(oneGB + 1024) {
		t.Error("checkDataCap should return false once already alerted")
	}

	// Reset alert
	resetDataCapAlert()
	if !checkDataCap(oneGB + 2048) {
		t.Error("checkDataCap should return true after alert reset")
	}

	// Cleanup
	setDataCap(0)
}
