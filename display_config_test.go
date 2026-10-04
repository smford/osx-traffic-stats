package main

import (
	"testing"
)

func TestDisplayConfig(t *testing.T) {
	// 1. Standard Bytes
	setUnitMode(UnitBytes)
	setStyleMode(StyleStandard)
	res := formatSpeedDynamic(1024*1024, true)
	if res != " 1.0 MB/s" {
		t.Errorf("expected ' 1.0 MB/s', got %q", res)
	}

	// 2. Compact Bytes
	setStyleMode(StyleCompact)
	res = formatSpeedDynamic(1024*1024, true)
	if res != " 1.0M" {
		t.Errorf("expected ' 1.0M', got %q", res)
	}

	// 3. Standard Bits
	setUnitMode(UnitBits)
	setStyleMode(StyleStandard)
	// 1 MB/s = 8 Mbps (8,000,000 bits) -> 8.4 Mbps or 8.0 Mbps depending on exact base
	res = formatSpeedDynamic(125000, true) // 125,000 bytes = 1,000,000 bits = 1.0 Mbps
	if res != " 1.0 Mbps" {
		t.Errorf("expected ' 1.0 Mbps', got %q", res)
	}

	// 4. Compact Bits
	setStyleMode(StyleCompact)
	res = formatSpeedDynamic(125000, true)
	if res != " 1.0M" {
		t.Errorf("expected ' 1.0M', got %q", res)
	}

	// Reset to defaults
	setUnitMode(UnitBytes)
	setStyleMode(StyleStandard)
}

func TestMenuBarIconMode(t *testing.T) {
	setMenuBarIconMode(MenuBarTextOnly)
	if getMenuBarIconMode() != MenuBarTextOnly {
		t.Errorf("expected MenuBarTextOnly")
	}

	setMenuBarIconMode(MenuBarGraphAndText)
	if getMenuBarIconMode() != MenuBarGraphAndText {
		t.Errorf("expected MenuBarGraphAndText")
	}

	setMenuBarIconMode(MenuBarGraphOnly)
	if getMenuBarIconMode() != MenuBarGraphOnly {
		t.Errorf("expected MenuBarGraphOnly")
	}

	setMenuBarIconMode(MenuBarTextOnly)
}
