package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"sync"

	"github.com/shirou/gopsutil/v3/net"
)

type InterfaceOption struct {
	Device      string
	DisplayName string
}

var (
	selectedInterfaceMu sync.RWMutex
	selectedInterface   string // empty string means "All Interfaces" (excluding loopback)
)

func getSelectedInterface() string {
	selectedInterfaceMu.RLock()
	defer selectedInterfaceMu.RUnlock()
	return selectedInterface
}

func setSelectedInterface(device string) {
	selectedInterfaceMu.Lock()
	selectedInterface = device
	selectedInterfaceMu.Unlock()
}

// getHardwarePortMap maps BSD device names (e.g. "en0") to friendly names (e.g. "Wi-Fi").
func getHardwarePortMap() map[string]string {
	portMap := make(map[string]string)
	cmd := exec.Command("networksetup", "-listallhardwareports")
	out, err := cmd.Output()
	if err != nil {
		return portMap
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	var currentPort string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Hardware Port:") {
			currentPort = strings.TrimSpace(strings.TrimPrefix(line, "Hardware Port:"))
		} else if strings.HasPrefix(line, "Device:") {
			device := strings.TrimSpace(strings.TrimPrefix(line, "Device:"))
			if device != "" && currentPort != "" {
				portMap[device] = currentPort
			}
			currentPort = ""
		}
	}
	return portMap
}

// listInterfaceOptions returns list of selectable interfaces.
func listInterfaceOptions() []InterfaceOption {
	ports := getHardwarePortMap()
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil
	}

	var options []InterfaceOption
	seen := make(map[string]bool)

	for _, io := range ioCounters {
		dev := io.Name
		if dev == "lo0" || seen[dev] {
			continue
		}
		seen[dev] = true

		name, hasName := ports[dev]
		if hasName {
			options = append(options, InterfaceOption{
				Device:      dev,
				DisplayName: name + " (" + dev + ")",
			})
		} else if strings.HasPrefix(dev, "utun") {
			options = append(options, InterfaceOption{
				Device:      dev,
				DisplayName: "VPN (" + dev + ")",
			})
		} else if strings.HasPrefix(dev, "en") {
			options = append(options, InterfaceOption{
				Device:      dev,
				DisplayName: "Ethernet/Wi-Fi (" + dev + ")",
			})
		}
	}

	return options
}

// fetchTrafficBytes queries the network stats for the currently selected interface (or sums all non-loopback interfaces).
func fetchTrafficBytes(targetDevice string) (sent, recv uint64, err error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return 0, 0, err
	}

	for _, c := range counters {
		if c.Name == "lo0" {
			continue // ignore local loopback
		}
		if targetDevice == "" || c.Name == targetDevice {
			sent += c.BytesSent
			recv += c.BytesRecv
		}
	}

	return sent, recv, nil
}
