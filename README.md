<p align="center">
  <img src="assets/logo.svg" alt="OSX Traffic Stats Logo" width="128" height="128" />
</p>

<h1 align="center">OSX Traffic Stats</h1>

<p align="center">
  A lightweight macOS menu bar utility that displays real-time network upload and download speeds next to the menu bar clock.
</p>

<p align="center">
  <a href="https://github.com/smford/osx-traffic-stats/releases"><img src="https://img.shields.io/github/v/release/smford/osx-traffic-stats?color=blue" alt="Latest Release" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT" /></a>
  <a href="https://github.com/smford/homebrew-tap"><img src="https://img.shields.io/badge/Homebrew-smford%2Ftap-orange.svg" alt="Homebrew" /></a>
  <a href="https://github.com/smford/osx-traffic-stats/actions/workflows/release.yml"><img src="https://img.shields.io/github/actions/workflow/status/smford/osx-traffic-stats/release.yml?branch=main" alt="Release Workflow" /></a>
</p>

---

## 📦 Installation

### Method 1: Homebrew Cask (Recommended — Native macOS App)

Install the native macOS `.app` bundle directly into `/Applications`:

```bash
brew install --cask smford/tap/osx-traffic-stats
```

Or add the tap first and then install:

```bash
brew tap smford/tap
brew install --cask osx-traffic-stats
```

#### Launching the App
- **Spotlight**: Press `Cmd + Space`, type `OSX Traffic Stats`, and press **Enter**.
- **Launchpad / Finder**: Click **OSX Traffic Stats** in `/Applications`.
- **Terminal**: Run `open -a OSXTrafficStats`.

#### Upgrading & Uninstalling
```bash
# Upgrade to the latest release
brew upgrade --cask osx-traffic-stats

# Uninstall
brew uninstall --cask osx-traffic-stats
```

---

### Method 2: Homebrew Formula (Standalone CLI Binary)

If you prefer only the standalone command-line binary installed to your `PATH` (`/opt/homebrew/bin`):

```bash
brew install --formula smford/tap/osx-traffic-stats
```

To run in the background from your terminal:
```bash
osx-traffic-stats &
```

---

### Method 3: Pre-built GitHub Releases

Download universal binaries (Apple Silicon & Intel) or the macOS application bundle (`OSXTrafficStats-v*.zip`) directly from [GitHub Releases](https://github.com/smford/osx-traffic-stats/releases):
1. Download `OSXTrafficStats-v*-macOS.zip`.
2. Extract the archive.
3. Drag `OSXTrafficStats.app` into your **Applications** folder.

---

### Method 4: Building from Source

**Prerequisites:** macOS (Darwin), Go 1.21+, and Xcode Command Line Tools (`xcode-select --install`).

```bash
git clone https://github.com/smford/osx-traffic-stats.git
cd osx-traffic-stats

# Build and install OSXTrafficStats.app directly to ~/Applications
make install

# Or build standalone binary in bin/osx-traffic-stats
make build

# Run directly from source
make run
```

---

## Features

- **Accessory Toolbar Behavior**: Operates seamlessly as a native macOS accessory/agent application (`NSApplicationActivationPolicyAccessory` & `LSUIElement`). It stays tucked in your menu bar with no Dock icon and no presence in the Command-Tab application switcher.
- **Real-Time Bandwidth Stats**: Displays current upload (`↑`) and download (`↓`) speeds updated every second.
- **Fixed-Width Menu Bar Layout**: Uses fixed status item length, monospaced tabular digits, and fixed-width speed slots so the menu bar item and neighboring icons never jump or shift left and right as network traffic fluctuates.
- **Adaptive Units**: Automatically scales transfer speeds (`KB/s`, `MB/s`, `GB/s`) and cumulative data (`B`, `KB`, `MB`, `GB`, `TB`).
- **Rich Menu Bar Dropdown**: Click the menu bar item to view:
  - Live upload and download speeds with dynamic in-menu sparkline graphs (`[ ▂▃▅█]`)
  - **Session Data Totals** and **Lifetime All-Time Totals**
  - One-click session stats reset
  - **Network Interface Selector** (All Physical Interfaces, Wi-Fi, Ethernet, VPN)
  - **Top Active Apps**: Submenu displaying real-time bandwidth consumption for the top network processes (via macOS `nettop`)
  - **Display Options**: Toggle speed units (Bytes/s vs. Bits/s), layout styles (Standard vs. Compact), and menu bar icon mode (Text Only, Dynamic Dual-Sparkline Graph + Text, or Graph Only)
  - **Configurable Refresh Rate (Battery Saver)**: Select polling frequencies: Fast (0.5s), Normal (1.0s), Battery Saver (2.0s), or Eco / Low Power (5.0s)
  - **Live Ping & Latency Indicator**: Measures round-trip ping time (`⚡ 6.7 ms`) with network quality grading (Excellent, Good, Fair, High Latency)
  - **Data Cap Alert**: Configurable session bandwidth warnings (1 GB, 2 GB, 5 GB, 10 GB, 20 GB) with native macOS system notifications
  - One-click toggle for the floating **Traffic Graph** window
  - **Launch at Login** toggle to start automatically on macOS boot
  - Native **About** dialog with copyright details, logo, and link to the GitHub repository
- **Persistent State**: Stats and preferences (interface selection, units, layout, and data cap limits) are automatically saved to `~/Library/Application Support/` and restored on launch.
- **Floating HUD Traffic Graph Window**: A sleek, translucent native macOS HUD window displaying:
  - Real-time 60-second dual-line bandwidth history with smooth curves and gradient area fills (Amber/Orange for Upload, Cyan/Blue for Download)
  - Live current and peak upload & download speed badges with one-click **Reset Peaks**
  - **Timeframe Selector**: Toggle between **1m**, **5m**, and **15m** history windows directly in the graph controls
  - **Top Active Apps Bar**: Real-time process bandwidth card displaying top bandwidth consumers right on the graph window
  - Auto-scaling dynamic Y-axis with dotted grid lines and time markers (`-60s`, `-30s`, `Now`)
  - **HUD Window Customization**: Adjust opacity presets (100% Solid, 85% Glass, 70% Subtle, 50% Translucent), toggle **Always on Top**, or snap to any screen corner (Top-Right, Bottom-Right, Top-Left, Bottom-Left)
  - Stays floating above windows for easy monitoring while working, draggable, resizable, and toggled from the menu bar
- **Hover Tooltip**: Detailed summary with live sparklines available on mouse hover.
- **Counter Reset Protection**: Handles network interface reconnections and counter rollovers cleanly without underflow spikes.
- **App Bundle Support**: Easily build and install as a native macOS application bundle (`OSXTrafficStats.app`).

---

## Development

A [`Makefile`](file:///Users/asc/git/osx-traffic-stats/Makefile) is provided for common development tasks:

| Target | Description |
|---|---|
| `make build` | Builds the standalone binary to `bin/osx-traffic-stats` with injected SemVer |
| `make app` | Builds the macOS application bundle `bin/OSXTrafficStats.app` with `LSUIElement` |
| `make package` | Builds universal (Apple Silicon & Intel) release `.zip` and `.tar.gz` |
| `make run` | Runs the application directly using `go run` |
| `make run-app` | Builds and launches `bin/OSXTrafficStats.app` using `open` |
| `make install` | Builds and installs `OSXTrafficStats.app` to `~/Applications` |
| `make test` | Runs unit tests |
| `make test-coverage` | Runs unit tests and produces an HTML coverage report |
| `make vet` | Runs `go vet` static analysis |
| `make clean` | Removes build and coverage artifacts |

---

## License

MIT License. See [LICENSE](LICENSE) for details.