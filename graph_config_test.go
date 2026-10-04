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
