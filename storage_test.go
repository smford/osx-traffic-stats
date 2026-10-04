package main

import (
	"encoding/json"
	"testing"
)

func TestAppStateSerialization(t *testing.T) {
	orig := &AppState{
		SessionSent:       1024 * 1024,
		SessionRecv:       50 * 1024 * 1024,
		TotalSent:         500 * 1024 * 1024,
		TotalRecv:         10 * 1024 * 1024 * 1024,
		SelectedInterface: "en0",
		UnitMode:          UnitBits,
		StyleMode:         StyleCompact,
		DataCapBytes:      5 * 1024 * 1024 * 1024,
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var loaded AppState
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if loaded.SessionSent != orig.SessionSent {
		t.Errorf("SessionSent = %d, want %d", loaded.SessionSent, orig.SessionSent)
	}
	if loaded.SessionRecv != orig.SessionRecv {
		t.Errorf("SessionRecv = %d, want %d", loaded.SessionRecv, orig.SessionRecv)
	}
	if loaded.TotalSent != orig.TotalSent {
		t.Errorf("TotalSent = %d, want %d", loaded.TotalSent, orig.TotalSent)
	}
	if loaded.TotalRecv != orig.TotalRecv {
		t.Errorf("TotalRecv = %d, want %d", loaded.TotalRecv, orig.TotalRecv)
	}
	if loaded.SelectedInterface != orig.SelectedInterface {
		t.Errorf("SelectedInterface = %q, want %q", loaded.SelectedInterface, orig.SelectedInterface)
	}
	if loaded.UnitMode != orig.UnitMode {
		t.Errorf("UnitMode = %v, want %v", loaded.UnitMode, orig.UnitMode)
	}
	if loaded.StyleMode != orig.StyleMode {
		t.Errorf("StyleMode = %v, want %v", loaded.StyleMode, orig.StyleMode)
	}
	if loaded.DataCapBytes != orig.DataCapBytes {
		t.Errorf("DataCapBytes = %d, want %d", loaded.DataCapBytes, orig.DataCapBytes)
	}
}
