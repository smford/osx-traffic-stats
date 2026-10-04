package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/getlantern/systray"
)

//go:embed assets/logo.svg
var appLogoSVG []byte

var (
	sessionSent atomic.Uint64
	sessionRecv atomic.Uint64
	totalSent   atomic.Uint64
	totalRecv   atomic.Uint64

	mUploadRate   *systray.MenuItem
	mDownloadRate *systray.MenuItem
	mSessionUp    *systray.MenuItem
	mSessionDown  *systray.MenuItem
	mTotalUp      *systray.MenuItem
	mTotalDown    *systray.MenuItem
	mToggleGraph  *systray.MenuItem

	appVersion = "1.0.0"

	refreshRateSignal = make(chan struct{}, 1)
)

const (
	appName      = "OSX Traffic Stats"
	appCopyright = "Copyright © 2026 smford"
	appLicense   = "MIT License"
	githubURL    = "https://github.com/smford/osx-traffic-stats"
)

func aboutMessage() string {
	return fmt.Sprintf("Version %s\n\n%s\nReleased under the %s\n\nGitHub:\n%s", appVersion, appCopyright, appLicense, githubURL)
}

func main() {
	showVersion := flag.Bool("version", false, "Print version and exit")
	showVersionShort := flag.Bool("v", false, "Print version and exit")
	flag.Parse()
	if *showVersion || *showVersionShort {
		fmt.Printf("%s %s\n", appName, appVersion)
		return
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	setAccessoryPolicy()
	configureFixedStatusItem()

	// Restore state from disk if available
	savedState, _ := loadAppState()
	if savedState != nil {
		sessionSent.Store(savedState.SessionSent)
		sessionRecv.Store(savedState.SessionRecv)
		totalSent.Store(savedState.TotalSent)
		totalRecv.Store(savedState.TotalRecv)
		if savedState.SelectedInterface != "" {
			setSelectedInterface(savedState.SelectedInterface)
		}
		setUnitMode(savedState.UnitMode)
		setStyleMode(savedState.StyleMode)
		setMenuBarIconMode(savedState.MenuBarIconMode)
		if savedState.RefreshRate > 0 {
			setRefreshRate(savedState.RefreshRate)
		}
		if savedState.GraphOpacity > 0 {
			setGraphOpacityConfig(savedState.GraphOpacity)
		}
		setGraphAlwaysOnTopConfig(savedState.GraphAlwaysOnTop)
		setDataCap(savedState.DataCapBytes)
	}

	systray.SetTitle("Loading...")
	systray.SetTooltip("Bandwidth Monitor")

	mUploadRate = systray.AddMenuItem("Upload: --", "Current upload speed")
	mUploadRate.Disable()
	mDownloadRate = systray.AddMenuItem("Download: --", "Current download speed")
	mDownloadRate.Disable()
	mLatency := systray.AddMenuItem("Ping: Measuring...", "Round-trip network latency to DNS gateway (1.1.1.1 / 8.8.8.8)")
	mLatency.Disable()

	StartLatencyMonitor(4*time.Second, func(status string) {
		mLatency.SetTitle(status)
	})

	systray.AddSeparator()

	mSessionUp = systray.AddMenuItem(fmt.Sprintf("Session Upload:   %s", formatBytes(sessionSent.Load())), "Total uploaded this session")
	mSessionUp.Disable()
	mSessionDown = systray.AddMenuItem(fmt.Sprintf("Session Download: %s", formatBytes(sessionRecv.Load())), "Total downloaded this session")
	mSessionDown.Disable()

	mTotalUp = systray.AddMenuItem(fmt.Sprintf("Lifetime Upload:  %s", formatBytes(totalSent.Load())), "All-time total uploaded")
	mTotalUp.Disable()
	mTotalDown = systray.AddMenuItem(fmt.Sprintf("Lifetime Download: %s", formatBytes(totalRecv.Load())), "All-time total downloaded")
	mTotalDown.Disable()

	mReset := systray.AddMenuItem("Reset Session Stats", "Reset session upload and download counters")
	mResetPeaks := systray.AddMenuItem("Reset Peak Speeds", "Reset peak upload and download speed records")

	systray.AddSeparator()

	mInterfaceMenu := systray.AddMenuItem("Network Interface", "Select network interface to monitor")
	mAllIfaces := mInterfaceMenu.AddSubMenuItemCheckbox("All Interfaces", "Monitor all non-loopback network interfaces", getSelectedInterface() == "")

	ifaceOptions := listInterfaceOptions()
	type ifaceItemPair struct {
		device string
		item   *systray.MenuItem
	}
	var ifaceItems []ifaceItemPair

	for _, opt := range ifaceOptions {
		subItem := mInterfaceMenu.AddSubMenuItemCheckbox(opt.DisplayName, "Monitor "+opt.Device, opt.Device == getSelectedInterface())
		ifaceItems = append(ifaceItems, ifaceItemPair{device: opt.Device, item: subItem})
	}

	systray.AddSeparator()

	mDisplayMenu := systray.AddMenuItem("Display Options", "Configure speed units and layout style")
	mUnitBytes := mDisplayMenu.AddSubMenuItemCheckbox("Units: Bytes (KB/s, MB/s)", "Display speed in bytes per second", getUnitMode() == UnitBytes)
	mUnitBits := mDisplayMenu.AddSubMenuItemCheckbox("Units: Bits (Kbps, Mbps)", "Display speed in bits per second", getUnitMode() == UnitBits)

	mStyleStandard := mDisplayMenu.AddSubMenuItemCheckbox("Style: Standard", "Full bandwidth labels", getStyleMode() == StyleStandard)
	mStyleCompact := mDisplayMenu.AddSubMenuItemCheckbox("Style: Compact (Notch-Friendly)", "Compact bandwidth labels", getStyleMode() == StyleCompact)

	mBarIconNone := mDisplayMenu.AddSubMenuItemCheckbox("Icon: Text Only", "Show bandwidth text only", getMenuBarIconMode() == MenuBarTextOnly)
	mBarIconBoth := mDisplayMenu.AddSubMenuItemCheckbox("Icon: Graph + Text", "Show real-time sparkline icon and text", getMenuBarIconMode() == MenuBarGraphAndText)
	mBarIconOnly := mDisplayMenu.AddSubMenuItemCheckbox("Icon: Graph Only", "Show real-time sparkline icon without text", getMenuBarIconMode() == MenuBarGraphOnly)

	mRefreshMenu := systray.AddMenuItem("Refresh Rate", "Configure polling frequency and battery saver")
	mRefreshFast := mRefreshMenu.AddSubMenuItemCheckbox("Fast (0.5s) - Ultra Smooth", "Update every 500ms", getRefreshRate() == RefreshFast)
	mRefreshNormal := mRefreshMenu.AddSubMenuItemCheckbox("Normal (1.0s) - Recommended", "Update every 1 second", getRefreshRate() == RefreshNormal)
	mRefreshBattery := mRefreshMenu.AddSubMenuItemCheckbox("Battery Saver (2.0s)", "Update every 2 seconds", getRefreshRate() == RefreshBattery)
	mRefreshEco := mRefreshMenu.AddSubMenuItemCheckbox("Eco / Low Power (5.0s)", "Update every 5 seconds", getRefreshRate() == RefreshEco)

	systray.AddSeparator()

	mDataCapMenu := systray.AddMenuItem("Data Cap Alert", "Set cellular/hotspot data usage warnings")
	type capItemPair struct {
		bytes uint64
		item  *systray.MenuItem
	}
	var capItems []capItemPair

	currentCap := getDataCap()
	for _, opt := range dataCapOptions {
		subItem := mDataCapMenu.AddSubMenuItemCheckbox(opt.Label, "Warn when session exceeds "+opt.Label, opt.Bytes == currentCap)
		capItems = append(capItems, capItemPair{bytes: opt.Bytes, item: subItem})
	}

	systray.AddSeparator()

	mTopProcsMenu := systray.AddMenuItem("Top Active Apps", "Applications currently consuming bandwidth")
	var procMenuItems []*systray.MenuItem
	for i := 0; i < 5; i++ {
		item := mTopProcsMenu.AddSubMenuItem("Monitoring...", "")
		item.Disable()
		procMenuItems = append(procMenuItems, item)
	}

	go func() {
		UpdateProcessTraffic()
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			procs := UpdateProcessTraffic()
			for i := 0; i < len(procMenuItems); i++ {
				if i < len(procs) {
					procMenuItems[i].SetTitle(FormatProcessItem(procs[i]))
					procMenuItems[i].Show()
				} else if i == 0 {
					procMenuItems[i].SetTitle("No active network traffic")
					procMenuItems[i].Show()
				} else {
					procMenuItems[i].Hide()
				}
			}
		}
	}()

	systray.AddSeparator()

	mToggleGraph = systray.AddMenuItem("Show Traffic Graph", "Open a real-time bandwidth graph window")

	mGraphOptions := systray.AddMenuItem("Graph Window Options", "Configure floating graph opacity, level, and position")
	mSnapTopRight := mGraphOptions.AddSubMenuItem("Snap: Top-Right Corner", "Position graph window in top-right corner")
	mSnapBottomRight := mGraphOptions.AddSubMenuItem("Snap: Bottom-Right Corner", "Position graph window in bottom-right corner")
	mSnapTopLeft := mGraphOptions.AddSubMenuItem("Snap: Top-Left Corner", "Position graph window in top-left corner")
	mSnapBottomLeft := mGraphOptions.AddSubMenuItem("Snap: Bottom-Left Corner", "Position graph window in bottom-left corner")

	mAlwaysOnTop := mGraphOptions.AddSubMenuItemCheckbox("Always on Top", "Keep floating graph above other windows", getGraphAlwaysOnTop())

	mOpacity100 := mGraphOptions.AddSubMenuItemCheckbox("Opacity: 100% (Solid)", "Full opacity", getGraphOpacity() >= 0.95)
	mOpacity85 := mGraphOptions.AddSubMenuItemCheckbox("Opacity: 85% (Glass)", "85% opacity", getGraphOpacity() >= 0.80 && getGraphOpacity() < 0.95)
	mOpacity70 := mGraphOptions.AddSubMenuItemCheckbox("Opacity: 70% (Subtle)", "70% opacity", getGraphOpacity() >= 0.60 && getGraphOpacity() < 0.80)
	mOpacity50 := mGraphOptions.AddSubMenuItemCheckbox("Opacity: 50% (Translucent)", "50% opacity", getGraphOpacity() < 0.60)

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
			resetDataCapAlert()
			mSessionUp.SetTitle("Session Upload:   0 B")
			mSessionDown.SetTitle("Session Download: 0 B")
			persistCurrentState()
		}
	}()

	for _, p := range capItems {
		pair := p
		go func() {
			for range pair.item.ClickedCh {
				setDataCap(pair.bytes)
				for _, other := range capItems {
					if other.bytes == pair.bytes {
						other.item.Check()
					} else {
						other.item.Uncheck()
					}
				}
				persistCurrentState()
			}
		}()
	}

	go func() {
		for range mResetPeaks.ClickedCh {
			resetPeaks()
		}
	}()

	go func() {
		for range mAllIfaces.ClickedCh {
			setSelectedInterface("")
			mAllIfaces.Check()
			for _, p := range ifaceItems {
				p.item.Uncheck()
			}
			persistCurrentState()
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
				persistCurrentState()
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
			persistCurrentState()
		}
	}()

	go func() {
		for range mUnitBits.ClickedCh {
			setUnitMode(UnitBits)
			mUnitBits.Check()
			mUnitBytes.Uncheck()
			persistCurrentState()
		}
	}()

	go func() {
		for range mStyleStandard.ClickedCh {
			setStyleMode(StyleStandard)
			mStyleStandard.Check()
			mStyleCompact.Uncheck()
			persistCurrentState()
		}
	}()

	go func() {
		for range mStyleCompact.ClickedCh {
			setStyleMode(StyleCompact)
			mStyleCompact.Check()
			mStyleStandard.Uncheck()
			persistCurrentState()
		}
	}()

	go func() {
		for range mBarIconNone.ClickedCh {
			setMenuBarIconMode(MenuBarTextOnly)
			mBarIconNone.Check()
			mBarIconBoth.Uncheck()
			mBarIconOnly.Uncheck()
			updateMenuBarGraph(0, 0, MenuBarTextOnly)
			persistCurrentState()
		}
	}()

	go func() {
		for range mBarIconBoth.ClickedCh {
			setMenuBarIconMode(MenuBarGraphAndText)
			mBarIconNone.Uncheck()
			mBarIconBoth.Check()
			mBarIconOnly.Uncheck()
			persistCurrentState()
		}
	}()

	go func() {
		for range mBarIconOnly.ClickedCh {
			setMenuBarIconMode(MenuBarGraphOnly)
			mBarIconNone.Uncheck()
			mBarIconBoth.Uncheck()
			mBarIconOnly.Check()
			systray.SetTitle("")
			persistCurrentState()
		}
	}()

	go func() {
		for range mRefreshFast.ClickedCh {
			setRefreshRate(RefreshFast)
			mRefreshFast.Check()
			mRefreshNormal.Uncheck()
			mRefreshBattery.Uncheck()
			mRefreshEco.Uncheck()
			select {
			case refreshRateSignal <- struct{}{}:
			default:
			}
			persistCurrentState()
		}
	}()

	go func() {
		for range mRefreshNormal.ClickedCh {
			setRefreshRate(RefreshNormal)
			mRefreshFast.Uncheck()
			mRefreshNormal.Check()
			mRefreshBattery.Uncheck()
			mRefreshEco.Uncheck()
			select {
			case refreshRateSignal <- struct{}{}:
			default:
			}
			persistCurrentState()
		}
	}()

	go func() {
		for range mRefreshBattery.ClickedCh {
			setRefreshRate(RefreshBattery)
			mRefreshFast.Uncheck()
			mRefreshNormal.Uncheck()
			mRefreshBattery.Check()
			mRefreshEco.Uncheck()
			select {
			case refreshRateSignal <- struct{}{}:
			default:
			}
			persistCurrentState()
		}
	}()

	go func() {
		for range mRefreshEco.ClickedCh {
			setRefreshRate(RefreshEco)
			mRefreshFast.Uncheck()
			mRefreshNormal.Uncheck()
			mRefreshBattery.Uncheck()
			mRefreshEco.Check()
			select {
			case refreshRateSignal <- struct{}{}:
			default:
			}
			persistCurrentState()
		}
	}()

	go func() {
		for range mSnapTopRight.ClickedCh {
			snapGraphWindow(0)
		}
	}()
	go func() {
		for range mSnapBottomRight.ClickedCh {
			snapGraphWindow(1)
		}
	}()
	go func() {
		for range mSnapTopLeft.ClickedCh {
			snapGraphWindow(2)
		}
	}()
	go func() {
		for range mSnapBottomLeft.ClickedCh {
			snapGraphWindow(3)
		}
	}()

	go func() {
		for range mAlwaysOnTop.ClickedCh {
			newVal := !getGraphAlwaysOnTop()
			setGraphAlwaysOnTopConfig(newVal)
			if newVal {
				mAlwaysOnTop.Check()
			} else {
				mAlwaysOnTop.Uncheck()
			}
			persistCurrentState()
		}
	}()

	go func() {
		for range mOpacity100.ClickedCh {
			setGraphOpacityConfig(1.0)
			mOpacity100.Check()
			mOpacity85.Uncheck()
			mOpacity70.Uncheck()
			mOpacity50.Uncheck()
			persistCurrentState()
		}
	}()
	go func() {
		for range mOpacity85.ClickedCh {
			setGraphOpacityConfig(0.85)
			mOpacity100.Uncheck()
			mOpacity85.Check()
			mOpacity70.Uncheck()
			mOpacity50.Uncheck()
			persistCurrentState()
		}
	}()
	go func() {
		for range mOpacity70.ClickedCh {
			setGraphOpacityConfig(0.70)
			mOpacity100.Uncheck()
			mOpacity85.Uncheck()
			mOpacity70.Check()
			mOpacity50.Uncheck()
			persistCurrentState()
		}
	}()
	go func() {
		for range mOpacity50.ClickedCh {
			setGraphOpacityConfig(0.50)
			mOpacity100.Uncheck()
			mOpacity85.Uncheck()
			mOpacity70.Uncheck()
			mOpacity50.Check()
			persistCurrentState()
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
	var tickCount int
	var lastTick time.Time
	const sparklineWidth = 10

	ticker := time.NewTicker(getRefreshDuration())
	defer ticker.Stop()

	for {
		select {
		case <-refreshRateSignal:
			ticker.Reset(getRefreshDuration())
			continue
		case now := <-ticker.C:
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
				var upDelta, downDelta uint64
				if currSent >= prevSent {
					upDelta = currSent - prevSent
				}
				if currRecv >= prevRecv {
					downDelta = currRecv - prevRecv
				}

				elapsed := 1.0
				if !lastTick.IsZero() {
					elapsed = now.Sub(lastTick).Seconds()
					if elapsed <= 0 {
						elapsed = 1.0
					}
				}

				upBytesPerSec := CalculateSpeedPerSecond(upDelta, elapsed)
				downBytesPerSec := CalculateSpeedPerSecond(downDelta, elapsed)

				totalUp := sessionSent.Add(upDelta)
				totalDown := sessionRecv.Add(downDelta)
				allTimeUp := totalSent.Add(upDelta)
				allTimeDown := totalRecv.Add(downDelta)
				checkDataCap(totalUp + totalDown)

				upSpeed := formatSpeedDynamic(upBytesPerSec, false)
				downSpeed := formatSpeedDynamic(downBytesPerSec, false)
				upSpeedFixed := formatSpeedDynamic(upBytesPerSec, true)
				downSpeedFixed := formatSpeedDynamic(downBytesPerSec, true)

				upHistory = append(upHistory, upBytesPerSec)
				if len(upHistory) > sparklineWidth {
					upHistory = upHistory[len(upHistory)-sparklineWidth:]
				}
				downHistory = append(downHistory, downBytesPerSec)
				if len(downHistory) > sparklineWidth {
					downHistory = downHistory[len(downHistory)-sparklineWidth:]
				}

				upSpark := renderSparkline(upHistory, sparklineWidth)
				downSpark := renderSparkline(downHistory, sparklineWidth)

				updateGraph(upBytesPerSec, downBytesPerSec)
				updateMenuBarGraph(upBytesPerSec, downBytesPerSec, getMenuBarIconMode())

				// Updates the text right next to the macOS clock with fixed-width layout
				if getMenuBarIconMode() == MenuBarGraphOnly {
					systray.SetTitle("")
				} else {
					systray.SetTitle(fmt.Sprintf("↑ %s  ↓ %s", upSpeedFixed, downSpeedFixed))
				}
				pingMs := currentLatencyMs.Load()
				var pingStr string
				if pingMs > 0 && pingMs < LatencyTimeoutMs {
					pingStr = fmt.Sprintf(" | Ping: %d ms", pingMs)
				}
				systray.SetTooltip(fmt.Sprintf("Bandwidth Monitor%s\n↑ %s [%s] (Session: %s, Total: %s)\n↓ %s [%s] (Session: %s, Total: %s)", pingStr, upSpeed, upSpark, formatBytes(totalUp), formatBytes(allTimeUp), downSpeed, downSpark, formatBytes(totalDown), formatBytes(allTimeDown)))

				mUploadRate.SetTitle(fmt.Sprintf("Upload:   ↑ %-8s  [%s]", upSpeed, upSpark))
				mDownloadRate.SetTitle(fmt.Sprintf("Download: ↓ %-8s  [%s]", downSpeed, downSpark))
				mSessionUp.SetTitle(fmt.Sprintf("Session Upload:   %s", formatBytes(totalUp)))
				mSessionDown.SetTitle(fmt.Sprintf("Session Download: %s", formatBytes(totalDown)))
				mTotalUp.SetTitle(fmt.Sprintf("Lifetime Upload:  %s", formatBytes(allTimeUp)))
				mTotalDown.SetTitle(fmt.Sprintf("Lifetime Download: %s", formatBytes(allTimeDown)))

				tickCount++
				if tickCount%15 == 0 {
					persistCurrentState()
				}
			} else {
				initialized = true
			}

			lastTick = now
			prevSent = currSent
			prevRecv = currRecv
		}
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

func onExit() {
	persistCurrentState()
}
