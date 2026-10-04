package main

import (
	"fmt"
	"sync/atomic"
)

type UnitMode int

const (
	UnitBytes UnitMode = iota // KB/s, MB/s (1024-based)
	UnitBits                  // Kbps, Mbps (1000-based bits)
)

type StyleMode int

const (
	StyleStandard StyleMode = iota // Standard fixed width
	StyleCompact                  // Compact mode for notches
)

var (
	currentUnitMode  atomic.Int32
	currentStyleMode atomic.Int32
)

func getUnitMode() UnitMode {
	return UnitMode(currentUnitMode.Load())
}

func setUnitMode(u UnitMode) {
	currentUnitMode.Store(int32(u))
}

func getStyleMode() StyleMode {
	return StyleMode(currentStyleMode.Load())
}

func setStyleMode(s StyleMode) {
	currentStyleMode.Store(int32(s))
}

// formatSpeedDynamic formats bandwidth based on selected unit and style
func formatSpeedDynamic(bytesPerSec uint64, fixed bool) string {
	if getUnitMode() == UnitBits {
		bitsPerSec := bytesPerSec * 8
		const (
			kb = 1000
			mb = 1000 * kb
			gb = 1000 * mb
		)
		if fixed {
			if getStyleMode() == StyleCompact {
				switch {
				case bitsPerSec >= 100*gb:
					return fmt.Sprintf("%4dG", bitsPerSec/gb)
				case bitsPerSec >= gb:
					return fmt.Sprintf("%4.1fG", float64(bitsPerSec)/float64(gb))
				case bitsPerSec >= 100*mb:
					return fmt.Sprintf("%4dM", bitsPerSec/mb)
				case bitsPerSec >= mb:
					return fmt.Sprintf("%4.1fM", float64(bitsPerSec)/float64(mb))
				default:
					return fmt.Sprintf("%4dK", bitsPerSec/kb)
				}
			}
			switch {
			case bitsPerSec >= 100*gb:
				return fmt.Sprintf("%4d Gbps", bitsPerSec/gb)
			case bitsPerSec >= gb:
				return fmt.Sprintf("%4.1f Gbps", float64(bitsPerSec)/float64(gb))
			case bitsPerSec >= 100*mb:
				return fmt.Sprintf("%4d Mbps", bitsPerSec/mb)
			case bitsPerSec >= mb:
				return fmt.Sprintf("%4.1f Mbps", float64(bitsPerSec)/float64(mb))
			default:
				return fmt.Sprintf("%4d Kbps", bitsPerSec/kb)
			}
		}
		// Non-fixed (e.g. for dropdown)
		switch {
		case bitsPerSec >= gb:
			return fmt.Sprintf("%.1f Gbps", float64(bitsPerSec)/float64(gb))
		case bitsPerSec >= mb:
			return fmt.Sprintf("%.1f Mbps", float64(bitsPerSec)/float64(mb))
		default:
			return fmt.Sprintf("%d Kbps", bitsPerSec/kb)
		}
	}

	// Bytes mode
	if fixed && getStyleMode() == StyleCompact {
		const (
			kb = 1024
			mb = 1024 * kb
			gb = 1024 * mb
		)
		switch {
		case bytesPerSec >= 100*gb:
			return fmt.Sprintf("%4dG", bytesPerSec/gb)
		case bytesPerSec >= gb:
			return fmt.Sprintf("%4.1fG", float64(bytesPerSec)/float64(gb))
		case bytesPerSec >= 100*mb:
			return fmt.Sprintf("%4dM", bytesPerSec/mb)
		case bytesPerSec >= mb:
			return fmt.Sprintf("%4.1fM", float64(bytesPerSec)/float64(mb))
		default:
			return fmt.Sprintf("%4dK", bytesPerSec/kb)
		}
	}

	if fixed {
		return formatSpeedFixed(bytesPerSec)
	}
	return formatSpeed(bytesPerSec)
}
