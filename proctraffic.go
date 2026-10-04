package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ProcessSnapshot stores cumulative traffic for a process.
type ProcessSnapshot struct {
	Name     string
	PID      int
	BytesIn  uint64
	BytesOut uint64
}

// ProcessTraffic represents calculated transfer rates for an active process.
type ProcessTraffic struct {
	Name     string
	PID      int
	RateIn   uint64 // bytes/sec
	RateOut  uint64 // bytes/sec
	TotalRate uint64 // RateIn + RateOut
}

var (
	procMutex      sync.RWMutex
	prevSnapshots  map[int]ProcessSnapshot
	lastSampleTime time.Time
	cachedTopProcs []ProcessTraffic
)

// ParseNettopOutput parses the CSV output from `nettop -P -x -L 1 -J bytes_in,bytes_out`.
func ParseNettopOutput(raw string) map[int]ProcessSnapshot {
	snapshots := make(map[int]ProcessSnapshot)
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ",bytes_in") || strings.HasPrefix(line, "time,") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}

		procIdentifier := parts[0]
		if procIdentifier == "" {
			continue
		}

		lastDot := strings.LastIndex(procIdentifier, ".")
		if lastDot == -1 {
			continue
		}

		name := procIdentifier[:lastDot]
		pidStr := procIdentifier[lastDot+1:]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		bIn, _ := strconv.ParseUint(parts[1], 10, 64)
		bOut, _ := strconv.ParseUint(parts[2], 10, 64)

		snapshots[pid] = ProcessSnapshot{
			Name:     name,
			PID:      pid,
			BytesIn:  bIn,
			BytesOut: bOut,
		}
	}

	return snapshots
}

// ComputeProcessRates calculates transfer rates between two snapshots.
func ComputeProcessRates(prev, curr map[int]ProcessSnapshot, intervalSec float64) []ProcessTraffic {
	if intervalSec <= 0 {
		intervalSec = 1.0
	}

	var results []ProcessTraffic

	for pid, c := range curr {
		p, exists := prev[pid]
		if !exists {
			continue
		}

		var rateIn, rateOut uint64
		if c.BytesIn >= p.BytesIn {
			rateIn = uint64(float64(c.BytesIn-p.BytesIn) / intervalSec)
		}
		if c.BytesOut >= p.BytesOut {
			rateOut = uint64(float64(c.BytesOut-p.BytesOut) / intervalSec)
		}

		total := rateIn + rateOut
		if total > 0 {
			results = append(results, ProcessTraffic{
				Name:      c.Name,
				PID:       pid,
				RateIn:    rateIn,
				RateOut:   rateOut,
				TotalRate: total,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalRate > results[j].TotalRate
	})

	return results
}

// FetchProcessSnapshots runs `nettop` once and parses the output.
func FetchProcessSnapshots() (map[int]ProcessSnapshot, error) {
	cmd := exec.Command("nettop", "-P", "-x", "-L", "1", "-J", "bytes_in,bytes_out")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return ParseNettopOutput(stdout.String()), nil
}

// UpdateProcessTraffic updates cached process bandwidth rates.
func UpdateProcessTraffic() []ProcessTraffic {
	curr, err := FetchProcessSnapshots()
	if err != nil {
		return GetCachedProcessTraffic()
	}

	now := time.Now()

	procMutex.Lock()
	defer procMutex.Unlock()

	if prevSnapshots == nil || lastSampleTime.IsZero() {
		prevSnapshots = curr
		lastSampleTime = now
		return cachedTopProcs
	}

	interval := now.Sub(lastSampleTime).Seconds()
	rates := ComputeProcessRates(prevSnapshots, curr, interval)
	prevSnapshots = curr
	lastSampleTime = now
	cachedTopProcs = rates

	return cachedTopProcs
}

// GetCachedProcessTraffic returns the latest computed process bandwidth list.
func GetCachedProcessTraffic() []ProcessTraffic {
	procMutex.RLock()
	defer procMutex.RUnlock()
	res := make([]ProcessTraffic, len(cachedTopProcs))
	copy(res, cachedTopProcs)
	return res
}

// FormatProcessItem returns a menu string for a process entry.
func FormatProcessItem(p ProcessTraffic) string {
	return fmt.Sprintf("%s: ↓ %s  ↑ %s", p.Name, formatSpeed(p.RateIn), formatSpeed(p.RateOut))
}

// FormatTopAppsSummary returns a concise summary of the top active apps for display on the graph window.
func FormatTopAppsSummary(maxApps int) string {
	procs := GetCachedProcessTraffic()
	if len(procs) == 0 {
		return "No active network traffic"
	}

	var parts []string
	count := 0
	for _, p := range procs {
		if p.TotalRate == 0 {
			continue
		}
		var speedStr string
		if p.RateIn > 0 && p.RateOut > 0 {
			speedStr = fmt.Sprintf("%s: ↓ %s  ↑ %s", p.Name, formatSpeed(p.RateIn), formatSpeed(p.RateOut))
		} else if p.RateIn > 0 {
			speedStr = fmt.Sprintf("%s: ↓ %s", p.Name, formatSpeed(p.RateIn))
		} else {
			speedStr = fmt.Sprintf("%s: ↑ %s", p.Name, formatSpeed(p.RateOut))
		}
		parts = append(parts, speedStr)
		count++
		if count >= maxApps {
			break
		}
	}

	if len(parts) == 0 {
		return "No active network traffic"
	}

	return strings.Join(parts, "   •   ")
}
