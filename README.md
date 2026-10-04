# osx-traffic-stats

A lightweight macOS menu bar utility that displays real-time network upload and download speeds next to the menu bar clock.

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
  - Native **About** dialog with copyright details and link to the GitHub repository
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

## Installation

### Via Homebrew (Recommended)

Install using the official Homebrew tap:

```bash
brew install smford/tap/osx-traffic-stats
```

Or add the tap first and then install:

```bash
brew tap smford/tap
brew install osx-traffic-stats
```

To upgrade:

```bash
brew update
brew upgrade osx-traffic-stats
```

### Pre-built Binaries & App Bundle

Download universal binaries (Apple Silicon & Intel) or the macOS application bundle (`OSXTrafficStats-v*.zip`) from [GitHub Releases](https://github.com/smford/osx-traffic-stats/releases).

## Prerequisites

- **macOS**: Built specifically for macOS (Darwin).
- **Go**: Version 1.21 or later.
- **Xcode Command Line Tools**: Required for Cgo support with macOS Cocoa APIs (`xcode-select --install`).

## Getting Started

### Building from Source

Clone the repository and build the binary or macOS `.app` bundle:

```bash
git clone https://github.com/smford/osx-traffic-stats.git
cd osx-traffic-stats

# Build standalone binary in bin/osx-traffic-stats
make build

# Or build native macOS .app bundle in bin/OSXTrafficStats.app
make app
```

### Running

To run directly from source:

```bash
make run
```

To run the standalone binary:

```bash
./bin/osx-traffic-stats
```

To launch the macOS `.app` bundle:

```bash
make run-app
```

To install the app bundle into `~/Applications`:

```bash
make install
```

To exit, click the menu bar item and select **Quit**.

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

## License

MIT