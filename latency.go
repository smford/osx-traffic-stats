package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	currentLatencyMs atomic.Uint64 // stored as uint64 milliseconds, or 999999 for offline
	latencyActive    atomic.Bool
)

const LatencyTimeoutMs uint64 = 999999

func init() {
	currentLatencyMs.Store(0)
}

// ParsePingOutput extracts the ping RTT in milliseconds from macOS ping output.
func ParsePingOutput(output string) (float64, error) {
	// Look for "round-trip min/avg/max/stddev = min/avg/max/stddev ms"
	if idx := strings.Index(output, "round-trip min/avg/max/stddev = "); idx != -1 {
		rest := output[idx+len("round-trip min/avg/max/stddev = "):]
		parts := strings.Split(rest, "/")
		if len(parts) >= 2 {
			avgStr := parts[1]
			val, err := strconv.ParseFloat(avgStr, 64)
			if err == nil {
				return val, nil
			}
		}
	}

	// Fallback to "time=X.XXX ms"
	if idx := strings.Index(output, "time="); idx != -1 {
		rest := output[idx+len("time="):]
		endIdx := strings.Index(rest, " ms")
		if endIdx != -1 {
			val, err := strconv.ParseFloat(strings.TrimSpace(rest[:endIdx]), 64)
			if err == nil {
				return val, nil
			}
		}
	}

	return 0, fmt.Errorf("could not parse latency from ping output")
}

// FormatLatencyStatus formats the latency into a human-readable status string.
func FormatLatencyStatus(ms float64, err error) string {
	if err != nil || ms < 0 || uint64(ms) >= LatencyTimeoutMs {
		return "Ping: Timeout / Offline"
	}

	switch {
	case ms < 30:
		return fmt.Sprintf("Ping: %.1f ms (Excellent)", ms)
	case ms < 80:
		return fmt.Sprintf("Ping: %.1f ms (Good)", ms)
	case ms < 150:
		return fmt.Sprintf("Ping: %.1f ms (Fair)", ms)
	default:
		return fmt.Sprintf("Ping: %.1f ms (High Latency)", ms)
	}
}

// MeasurePingLatency runs a single ICMP ping against a target host with a 1-second timeout.
func MeasurePingLatency(target string) (float64, error) {
	if target == "" {
		target = "1.1.1.1"
	}
	cmd := exec.Command("/sbin/ping", "-c", "1", "-t", "1", target)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return 0, err
	}
	return ParsePingOutput(stdout.String())
}

// StartLatencyMonitor launches a background loop that measures latency every interval.
func StartLatencyMonitor(interval time.Duration, updateFn func(string)) {
	if latencyActive.Swap(true) {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Run once immediately
		ms, err := MeasurePingLatency("1.1.1.1")
		if err != nil {
			// Fallback to Google DNS 8.8.8.8
			ms, err = MeasurePingLatency("8.8.8.8")
		}
		if err != nil {
			currentLatencyMs.Store(LatencyTimeoutMs)
		} else {
			currentLatencyMs.Store(uint64(ms))
		}
		if updateFn != nil {
			updateFn(FormatLatencyStatus(ms, err))
		}

		for range ticker.C {
			ms, err := MeasurePingLatency("1.1.1.1")
			if err != nil {
				ms, err = MeasurePingLatency("8.8.8.8")
			}

			if err != nil {
				currentLatencyMs.Store(LatencyTimeoutMs)
			} else {
				currentLatencyMs.Store(uint64(ms))
			}

			if updateFn != nil {
				updateFn(FormatLatencyStatus(ms, err))
			}
		}
	}()
}
