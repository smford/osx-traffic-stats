package main

import "strings"

var sparkLevels = []rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// renderSparkline returns a sparkline string of the given width representing the recent values.
// Values are scaled relative to the maximum non-zero value in the slice.
func renderSparkline(values []uint64, width int) string {
	if width <= 0 {
		return ""
	}

	// Pad or truncate to exact width
	data := make([]uint64, width)
	if len(values) >= width {
		copy(data, values[len(values)-width:])
	} else {
		offset := width - len(values)
		copy(data[offset:], values)
	}

	var maxVal uint64
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	var sb strings.Builder
	sb.Grow(width * 4)

	for _, v := range data {
		if maxVal == 0 || v == 0 {
			sb.WriteRune(sparkLevels[0])
			continue
		}

		level := int((v * 7) / maxVal)
		if level == 0 {
			level = 1
		}
		if level >= len(sparkLevels) {
			level = len(sparkLevels) - 1
		}
		sb.WriteRune(sparkLevels[level])
	}

	return sb.String()
}
