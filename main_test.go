package main

import (
	"strings"
	"testing"
)

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "Zero bytes",
			bytes:    0,
			expected: "0 KB/s",
		},
		{
			name:     "Under 1 KB",
			bytes:    512,
			expected: "0 KB/s",
		},
		{
			name:     "Exact 1 KB",
			bytes:    1024,
			expected: "1 KB/s",
		},
		{
			name:     "Mid KB",
			bytes:    150 * 1024,
			expected: "150 KB/s",
		},
		{
			name:     "Exact 1 MB",
			bytes:    1024 * 1024,
			expected: "1.0 MB/s",
		},
		{
			name:     "Fractional MB",
			bytes:    2560 * 1024, // 2.5 MB
			expected: "2.5 MB/s",
		},
		{
			name:     "Exact 1 GB",
			bytes:    1024 * 1024 * 1024,
			expected: "1.0 GB/s",
		},
		{
			name:     "Fractional GB",
			bytes:    uint64(1.5 * 1024 * 1024 * 1024),
			expected: "1.5 GB/s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSpeed(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatSpeed(%d) = %q; want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatSpeedFixed(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{name: "Zero bytes", bytes: 0, expected: "   0 KB/s"},
		{name: "Under 1 KB", bytes: 512, expected: "   0 KB/s"},
		{name: "Exact 1 KB", bytes: 1024, expected: "   1 KB/s"},
		{name: "Mid KB", bytes: 150 * 1024, expected: " 150 KB/s"},
		{name: "Exact 1 MB", bytes: 1024 * 1024, expected: " 1.0 MB/s"},
		{name: "Fractional MB", bytes: 2560 * 1024, expected: " 2.5 MB/s"},
		{name: "Large MB (>100)", bytes: 150 * 1024 * 1024, expected: " 150 MB/s"},
		{name: "Exact 1 GB", bytes: 1024 * 1024 * 1024, expected: " 1.0 GB/s"},
		{name: "Fractional GB", bytes: uint64(1.5 * 1024 * 1024 * 1024), expected: " 1.5 GB/s"},
		{name: "Large GB (>100)", bytes: 120 * 1024 * 1024 * 1024, expected: " 120 GB/s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSpeedFixed(tt.bytes)
			if len(result) != 9 {
				t.Errorf("formatSpeedFixed(%d) length = %d; want 9 (result: %q)", tt.bytes, len(result), result)
			}
			if result != tt.expected {
				t.Errorf("formatSpeedFixed(%d) = %q; want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "Zero bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "Bytes below 1 KB",
			bytes:    500,
			expected: "500 B",
		},
		{
			name:     "Exact 1 KB",
			bytes:    1024,
			expected: "1.0 KB",
		},
		{
			name:     "Fractional KB",
			bytes:    1536,
			expected: "1.5 KB",
		},
		{
			name:     "Exact 1 MB",
			bytes:    1024 * 1024,
			expected: "1.0 MB",
		},
		{
			name:     "Fractional MB",
			bytes:    int64ToUint64(2.5 * 1024 * 1024),
			expected: "2.5 MB",
		},
		{
			name:     "Exact 1 GB",
			bytes:    1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "Fractional GB",
			bytes:    int64ToUint64(3.75 * 1024 * 1024 * 1024),
			expected: "3.75 GB",
		},
		{
			name:     "Exact 1 TB",
			bytes:    1024 * 1024 * 1024 * 1024,
			expected: "1.00 TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q; want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func int64ToUint64(f float64) uint64 {
	return uint64(f)
}

func TestAboutMessage(t *testing.T) {
	msg := aboutMessage()

	expectedSubstrings := []string{
		appVersion,
		appCopyright,
		appLicense,
		githubURL,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(msg, sub) {
			t.Errorf("aboutMessage() missing expected substring %q; got %q", sub, msg)
		}
	}
}
