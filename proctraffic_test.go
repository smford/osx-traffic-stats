package main

import (
	"testing"
)

func TestParseNettopOutput(t *testing.T) {
	raw := `,bytes_in,bytes_out,
syslogd.345,100,500,
mDNSResponder.517,2000,4000,
Google Chrome H.17889,1048576,2097152,
invalid_line_without_pid,0,0
,0,0
`
	snapshots := ParseNettopOutput(raw)
	if len(snapshots) != 3 {
		t.Fatalf("expected 3 parsed processes, got %d", len(snapshots))
	}

	chrome, ok := snapshots[17889]
	if !ok {
		t.Fatalf("expected PID 17889 to be present")
	}
	if chrome.Name != "Google Chrome H" {
		t.Errorf("expected name 'Google Chrome H', got '%s'", chrome.Name)
	}
	if chrome.BytesIn != 1048576 || chrome.BytesOut != 2097152 {
		t.Errorf("unexpected bytes in/out: %d / %d", chrome.BytesIn, chrome.BytesOut)
	}
}

func TestComputeProcessRates(t *testing.T) {
	prev := map[int]ProcessSnapshot{
		100: {Name: "curl", PID: 100, BytesIn: 1000, BytesOut: 500},
		200: {Name: "zoom", PID: 200, BytesIn: 5000, BytesOut: 10000},
	}
	curr := map[int]ProcessSnapshot{
		100: {Name: "curl", PID: 100, BytesIn: 2000, BytesOut: 500},   // +1000 in, 0 out
		200: {Name: "zoom", PID: 200, BytesIn: 15000, BytesOut: 30000}, // +10000 in, +20000 out
	}

	rates := ComputeProcessRates(prev, curr, 1.0)
	if len(rates) != 2 {
		t.Fatalf("expected 2 active processes, got %d", len(rates))
	}

	// Zoom should be first due to higher total rate
	if rates[0].Name != "zoom" {
		t.Errorf("expected rates[0] to be zoom, got %s", rates[0].Name)
	}
	if rates[0].RateIn != 10000 || rates[0].RateOut != 20000 {
		t.Errorf("unexpected zoom rates: in=%d, out=%d", rates[0].RateIn, rates[0].RateOut)
	}

	if rates[1].Name != "curl" {
		t.Errorf("expected rates[1] to be curl, got %s", rates[1].Name)
	}
	if rates[1].RateIn != 1000 || rates[1].RateOut != 0 {
		t.Errorf("unexpected curl rates: in=%d, out=%d", rates[1].RateIn, rates[1].RateOut)
	}
}

func TestFormatProcessItem(t *testing.T) {
	proc := ProcessTraffic{
		Name:    "Safari",
		PID:     1234,
		RateIn:  1024 * 1024,
		RateOut: 500 * 1024,
	}
	item := FormatProcessItem(proc)
	if item != "Safari: ↓ 1.0 MB/s  ↑ 500 KB/s" {
		t.Errorf("unexpected format: %s", item)
	}
}

func TestFormatTopAppsSummary(t *testing.T) {
	// Empty case
	procMutex.Lock()
	cachedTopProcs = nil
	procMutex.Unlock()

	summary := FormatTopAppsSummary(3)
	if summary != "No active network traffic" {
		t.Errorf("expected 'No active network traffic', got %q", summary)
	}

	// Populated case
	procMutex.Lock()
	cachedTopProcs = []ProcessTraffic{
		{Name: "Google Chrome", PID: 101, RateIn: 1024 * 1024, RateOut: 128 * 1024, TotalRate: 1024*1024 + 128*1024},
		{Name: "curl", PID: 102, RateIn: 500 * 1024, RateOut: 0, TotalRate: 500 * 1024},
	}
	procMutex.Unlock()

	summary = FormatTopAppsSummary(3)
	expected := "Google Chrome: ↓ 1.0 MB/s  ↑ 128 KB/s   •   curl: ↓ 500 KB/s"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
