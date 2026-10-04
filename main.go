package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/getlantern/systray"
)

var (
	sessionSent atomic.Uint64
	sessionRecv atomic.Uint64

	mUploadRate   *systray.MenuItem
	mDownloadRate *systray.MenuItem
	mSessionUp    *systray.MenuItem
	mSessionDown  *systray.MenuItem
	mToggleGraph  *systray.MenuItem
)

const (
	appName      = "OSX Traffic Stats"
	appVersion   = "1.0.0"
	appCopyright = "Copyright © 2026 smford"
	appLicense   = "MIT License"
	githubURL    = "https://github.com/smford/osx-traffic-stats"
)

func aboutMessage() string {
	return fmt.Sprintf("Version %s\n\n%s\nReleased under the %s\n\nGitHub:\n%s", appVersion, appCopyright, appLicense, githubURL)
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	setAccessoryPolicy()
	configureFixedStatusItem()

	systray.SetTitle("Loading...")
	systray.SetTooltip("Bandwidth Monitor")

	mUploadRate = systray.AddMenuItem("Upload: --", "Current upload speed")
	mUploadRate.Disable()
	mDownloadRate = systray.AddMenuItem("Download: --", "Current download speed")
	mDownloadRate.Disable()

	systray.AddSeparator()

	mSessionUp = systray.AddMenuItem("Session Upload: 0 B", "Total uploaded this session")
	mSessionUp.Disable()
	mSessionDown = systray.AddMenuItem("Session Download: 0 B", "Total downloaded this session")
	mSessionDown.Disable()

	mReset := systray.AddMenuItem("Reset Session Stats", "Reset session upload and download counters")

	systray.AddSeparator()

	mInterfaceMenu := systray.AddMenuItem("Network Interface", "Select network interface to monitor")
	mAllIfaces := mInterfaceMenu.AddSubMenuItemCheckbox("All Interfaces", "Monitor all non-loopback network interfaces", true)

	ifaceOptions := listInterfaceOptions()
	type ifaceItemPair struct {
		device string
		item   *systray.MenuItem
	}
	var ifaceItems []ifaceItemPair

	for _, opt := range ifaceOptions {
		subItem := mInterfaceMenu.AddSubMenuItemCheckbox(opt.DisplayName, "Monitor "+opt.Device, false)
		ifaceItems = append(ifaceItems, ifaceItemPair{device: opt.Device, item: subItem})
	}

	systray.AddSeparator()

	mDisplayMenu := systray.AddMenuItem("Display Options", "Configure speed units and layout style")
	mUnitBytes := mDisplayMenu.AddSubMenuItemCheckbox("Units: Bytes (KB/s, MB/s)", "Display speed in bytes per second", true)
	mUnitBits := mDisplayMenu.AddSubMenuItemCheckbox("Units: Bits (Kbps, Mbps)", "Display speed in bits per second", false)

	mStyleStandard := mDisplayMenu.AddSubMenuItemCheckbox("Style: Standard", "Full bandwidth labels", true)
	mStyleCompact := mDisplayMenu.AddSubMenuItemCheckbox("Style: Compact (Notch-Friendly)", "Compact bandwidth labels", false)

	systray.AddSeparator()

	mToggleGraph = systray.AddMenuItem("Show Traffic Graph", "Open a real-time bandwidth graph window")

	systray.AddSeparator()

	mLaunchAtLogin := systray.AddMenuItemCheckbox("Launch at Login", "Start Bandwidth Monitor automatically on login", isLaunchAtLoginEnabled())

	systray.AddSeparator()

	mAbout := systray.AddMenuItem("About OSX Traffic Stats", "Show application and copyright information")

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit Bandwidth Monitor")

	go func() {
		for range mReset.ClickedCh {
			sessionSent.Store(0)
			sessionRecv.Store(0)
			mSessionUp.SetTitle("Session Upload: 0 B")
			mSessionDown.SetTitle("Session Download: 0 B")
		}
	}()

	go func() {
		for range mAllIfaces.ClickedCh {
			setSelectedInterface("")
			mAllIfaces.Check()
			for _, p := range ifaceItems {
				p.item.Uncheck()
			}
		}
	}()

	for _, p := range ifaceItems {
		pair := p
		go func() {
			for range pair.item.ClickedCh {
				setSelectedInterface(pair.device)
				mAllIfaces.Uncheck()
				for _, other := range ifaceItems {
					if other.device == pair.device {
						other.item.Check()
					} else {
						other.item.Uncheck()
					}
				}
			}
		}()
	}

	go func() {
		for range mToggleGraph.ClickedCh {
			toggleTrafficGraph()
		}
	}()

	go func() {
		for range mLaunchAtLogin.ClickedCh {
			newVal := !mLaunchAtLogin.Checked()
			if err := setLaunchAtLogin(newVal); err == nil {
				if newVal {
					mLaunchAtLogin.Check()
				} else {
					mLaunchAtLogin.Uncheck()
				}
			}
		}
	}()

	go func() {
		for range mUnitBytes.ClickedCh {
			setUnitMode(UnitBytes)
			mUnitBytes.Check()
			mUnitBits.Uncheck()
		}
	}()

	go func() {
		for range mUnitBits.ClickedCh {
			setUnitMode(UnitBits)
			mUnitBits.Check()
			mUnitBytes.Uncheck()
		}
	}()

	go func() {
		for range mStyleStandard.ClickedCh {
			setStyleMode(StyleStandard)
			mStyleStandard.Check()
			mStyleCompact.Uncheck()
		}
	}()

	go func() {
		for range mStyleCompact.ClickedCh {
			setStyleMode(StyleCompact)
			mStyleCompact.Check()
			mStyleStandard.Uncheck()
		}
	}()

	go func() {
		for range mAbout.ClickedCh {
			showAboutBox()
		}
	}()

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go monitorTraffic()
}

func monitorTraffic() {
	var prevSent, prevRecv uint64
	var initialized bool
	var lastInterface string
	var upHistory, downHistory []uint64
	const sparklineWidth = 10

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		currentIface := getSelectedInterface()
		if currentIface != lastInterface {
			initialized = false
			lastInterface = currentIface
		}

		currSent, currRecv, err := fetchTrafficBytes(currentIface)
		if err != nil {
			continue
		}

		if initialized {
			var upBytes, downBytes uint64
			if currSent >= prevSent {
				upBytes = currSent - prevSent
			}
			if currRecv >= prevRecv {
				downBytes = currRecv - prevRecv
			}

			totalUp := sessionSent.Add(upBytes)
			totalDown := sessionRecv.Add(downBytes)

			upSpeed := formatSpeedDynamic(upBytes, false)
			downSpeed := formatSpeedDynamic(downBytes, false)
			upSpeedFixed := formatSpeedDynamic(upBytes, true)
			downSpeedFixed := formatSpeedDynamic(downBytes, true)

			upHistory = append(upHistory, upBytes)
			if len(upHistory) > sparklineWidth {
				upHistory = upHistory[len(upHistory)-sparklineWidth:]
			}
			downHistory = append(downHistory, downBytes)
			if len(downHistory) > sparklineWidth {
				downHistory = downHistory[len(downHistory)-sparklineWidth:]
			}

			upSpark := renderSparkline(upHistory, sparklineWidth)
			downSpark := renderSparkline(downHistory, sparklineWidth)

			updateGraph(upBytes, downBytes)

			// Updates the text right next to the macOS clock with fixed-width layout
			systray.SetTitle(fmt.Sprintf("↑ %s  ↓ %s", upSpeedFixed, downSpeedFixed))
			systray.SetTooltip(fmt.Sprintf("Bandwidth Monitor\n↑ %s [%s] (Total: %s)\n↓ %s [%s] (Total: %s)", upSpeed, upSpark, formatBytes(totalUp), downSpeed, downSpark, formatBytes(totalDown)))

			mUploadRate.SetTitle(fmt.Sprintf("Upload:   ↑ %-8s  [%s]", upSpeed, upSpark))
			mDownloadRate.SetTitle(fmt.Sprintf("Download: ↓ %-8s  [%s]", downSpeed, downSpark))
			mSessionUp.SetTitle(fmt.Sprintf("Session Upload:   %s", formatBytes(totalUp)))
			mSessionDown.SetTitle(fmt.Sprintf("Session Download: %s", formatBytes(totalDown)))
		} else {
			initialized = true
		}

		prevSent = currSent
		prevRecv = currRecv
	}
}

func formatSpeedFixed(bytesPerSec uint64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case bytesPerSec >= 100*gb:
		return fmt.Sprintf("%4d GB/s", bytesPerSec/gb)
	case bytesPerSec >= gb:
		return fmt.Sprintf("%4.1f GB/s", float64(bytesPerSec)/float64(gb))
	case bytesPerSec >= 100*mb:
		return fmt.Sprintf("%4d MB/s", bytesPerSec/mb)
	case bytesPerSec >= mb:
		return fmt.Sprintf("%4.1f MB/s", float64(bytesPerSec)/float64(mb))
	default:
		return fmt.Sprintf("%4d KB/s", bytesPerSec/kb)
	}
}

func formatSpeed(bytesPerSec uint64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case bytesPerSec >= gb:
		return fmt.Sprintf("%.1f GB/s", float64(bytesPerSec)/float64(gb))
	case bytesPerSec >= mb:
		return fmt.Sprintf("%.1f MB/s", float64(bytesPerSec)/float64(mb))
	default:
		return fmt.Sprintf("%d KB/s", bytesPerSec/kb)
	}
}

func formatBytes(b uint64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
		tb = 1024 * gb
	)

	switch {
	case b >= tb:
		return fmt.Sprintf("%.2f TB", float64(b)/float64(tb))
	case b >= gb:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(gb))
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func onExit() {}
