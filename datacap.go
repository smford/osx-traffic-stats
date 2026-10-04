package main

import (
	"fmt"
	"os/exec"
	"sync/atomic"
)

type DataCapOption struct {
	Label string
	Bytes uint64
}

var dataCapOptions = []DataCapOption{
	{Label: "Disabled", Bytes: 0},
	{Label: "1 GB", Bytes: 1 * 1024 * 1024 * 1024},
	{Label: "2 GB", Bytes: 2 * 1024 * 1024 * 1024},
	{Label: "5 GB", Bytes: 5 * 1024 * 1024 * 1024},
	{Label: "10 GB", Bytes: 10 * 1024 * 1024 * 1024},
	{Label: "20 GB", Bytes: 20 * 1024 * 1024 * 1024},
}

var (
	selectedDataCapBytes atomic.Uint64
	dataCapAlerted       atomic.Bool
)

func getDataCap() uint64 {
	return selectedDataCapBytes.Load()
}

func setDataCap(bytes uint64) {
	selectedDataCapBytes.Store(bytes)
	dataCapAlerted.Store(false)
}

func checkDataCap(totalSessionBytes uint64) bool {
	cap := getDataCap()
	if cap == 0 {
		return false
	}
	if totalSessionBytes >= cap && !dataCapAlerted.Load() {
		dataCapAlerted.Store(true)
		sendDataCapNotification(totalSessionBytes, cap)
		return true
	}
	return false
}

func resetDataCapAlert() {
	dataCapAlerted.Store(false)
}

func sendDataCapNotification(currentBytes, capBytes uint64) {
	title := "OSX Traffic Stats: Data Alert"
	msg := fmt.Sprintf("Session traffic reached %s (limit: %s)!", formatBytes(currentBytes), formatBytes(capBytes))
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "Ping"`, msg, title)
	_ = exec.Command("osascript", "-e", script).Start()
}
