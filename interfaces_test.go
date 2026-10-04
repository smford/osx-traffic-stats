package main

import (
	"testing"
)

func TestInterfaceSelection(t *testing.T) {
	setSelectedInterface("en0")
	if got := getSelectedInterface(); got != "en0" {
		t.Errorf("getSelectedInterface() = %q; want %q", got, "en0")
	}

	setSelectedInterface("")
	if got := getSelectedInterface(); got != "" {
		t.Errorf("getSelectedInterface() = %q; want %q", got, "")
	}
}

func TestFetchTrafficBytes(t *testing.T) {
	sent, recv, err := fetchTrafficBytes("")
	if err != nil {
		t.Fatalf("fetchTrafficBytes(\"\") error: %v", err)
	}
	// At least 0 bytes should be returned without error
	t.Logf("Total physical traffic: sent=%d, recv=%d", sent, recv)
}
