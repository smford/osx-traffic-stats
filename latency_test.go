package main

import (
	"errors"
	"testing"
)

func TestParsePingOutput(t *testing.T) {
	// macOS round-trip summary line
	sample1 := `PING 1.1.1.1 (1.1.1.1): 56 data bytes
64 bytes from 1.1.1.1: icmp_seq=0 ttl=57 time=8.450 ms

--- 1.1.1.1 ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 8.450/8.450/8.450/nan ms
`
	ms, err := ParsePingOutput(sample1)
	if err != nil {
		t.Fatalf("unexpected error parsing sample1: %v", err)
	}
	if ms != 8.450 {
		t.Errorf("expected 8.450 ms, got %f", ms)
	}

	// Line with time= only
	sample2 := `64 bytes from 8.8.8.8: icmp_seq=0 ttl=115 time=23.120 ms`
	ms2, err2 := ParsePingOutput(sample2)
	if err2 != nil {
		t.Fatalf("unexpected error parsing sample2: %v", err2)
	}
	if ms2 != 23.120 {
		t.Errorf("expected 23.120 ms, got %f", ms2)
	}

	// Invalid output
	_, err3 := ParsePingOutput("Request timeout for icmp_seq 0")
	if err3 == nil {
		t.Errorf("expected error for timed out ping")
	}
}

func TestFormatLatencyStatus(t *testing.T) {
	tests := []struct {
		ms       float64
		err      error
		expected string
	}{
		{ms: 12.4, err: nil, expected: "Ping: 12.4 ms (Excellent)"},
		{ms: 45.0, err: nil, expected: "Ping: 45.0 ms (Good)"},
		{ms: 110.5, err: nil, expected: "Ping: 110.5 ms (Fair)"},
		{ms: 220.0, err: nil, expected: "Ping: 220.0 ms (High Latency)"},
		{ms: 0, err: errors.New("timeout"), expected: "Ping: Timeout / Offline"},
		{ms: float64(LatencyTimeoutMs), err: nil, expected: "Ping: Timeout / Offline"},
	}

	for _, tt := range tests {
		got := FormatLatencyStatus(tt.ms, tt.err)
		if got != tt.expected {
			t.Errorf("FormatLatencyStatus(%f, %v) = %q; want %q", tt.ms, tt.err, got, tt.expected)
		}
	}
}
