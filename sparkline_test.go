package main

import (
	"testing"
)

func TestRenderSparkline(t *testing.T) {
	tests := []struct {
		name     string
		values   []uint64
		width    int
		expected string
	}{
		{
			name:     "Zero width",
			values:   []uint64{10, 20},
			width:    0,
			expected: "",
		},
		{
			name:     "Negative width",
			values:   []uint64{10, 20},
			width:    -1,
			expected: "",
		},
		{
			name:     "All zeroes",
			values:   []uint64{0, 0, 0, 0, 0},
			width:    5,
			expected: "     ",
		},
		{
			name:     "Empty values padded to width",
			values:   []uint64{},
			width:    4,
			expected: "    ",
		},
		{
			name:     "Ascending values",
			values:   []uint64{0, 20, 40, 60, 80, 100, 120},
			width:    7,
			expected: " ▂▃▄▅▆█",
		},
		{
			name:     "Truncated to width",
			values:   []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			width:    3,
			expected: "▆▇█",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderSparkline(tt.values, tt.width)
			if got != tt.expected {
				t.Errorf("renderSparkline(%v, %d) = %q; want %q", tt.values, tt.width, got, tt.expected)
			}
		})
	}
}
