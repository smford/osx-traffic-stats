# osx-traffic-stats

A lightweight macOS menu bar utility that displays real-time network upload and download speeds next to the menu bar clock.

## Features

- **Accessory Toolbar Behavior**: Operates seamlessly as a native macOS accessory/agent application (`NSApplicationActivationPolicyAccessory` & `LSUIElement`). It stays tucked in your menu bar with no Dock icon and no presence in the Command-Tab application switcher.
- **Real-Time Bandwidth Stats**: Displays current upload (`↑`) and download (`↓`) speeds updated every second.
- **Fixed-Width Menu Bar Layout**: Uses fixed status item length, monospaced tabular digits, and fixed-width speed slots so the menu bar item and neighboring icons never jump or shift left and right as network traffic fluctuates.
- **Adaptive Units**: Automatically scales transfer speeds (`KB/s`, `MB/s`, `GB/s`) and cumulative data (`B`, `KB`, `MB`, `GB`, `TB`).
- **Rich Menu Bar Dropdown**: Click the menu bar item to view:
  - Live upload and download speeds with dynamic in-menu sparkline graphs (`[ ▂▃▅█]`)
  - Session data totals (cumulative uploaded and downloaded bytes)
  - One-click session stats reset
  - One-click toggle for the floating **Traffic Graph** window
  - Native **About** dialog with copyright details and link to the GitHub repository
- **Floating HUD Traffic Graph Window**: A sleek, translucent native macOS HUD window displaying:
  - Real-time 60-second dual-line bandwidth history with smooth curves and gradient area fills (Amber/Orange for Upload, Cyan/Blue for Download)
  - Live current and peak upload & download speed badges
  - Auto-scaling dynamic Y-axis with dotted grid lines and time markers (`-60s`, `-30s`, `Now`)
  - Stays floating above windows for easy monitoring while working, draggable, resizable, and toggled from the menu bar
- **Hover Tooltip**: Detailed summary with live sparklines available on mouse hover.
- **Counter Reset Protection**: Handles network interface reconnections and counter rollovers cleanly without underflow spikes.
- **App Bundle Support**: Easily build and install as a native macOS application bundle (`OSXTrafficStats.app`).

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
| `make build` | Builds the standalone binary to `bin/osx-traffic-stats` |
| `make app` | Builds the macOS application bundle `bin/OSXTrafficStats.app` with `LSUIElement` |
| `make run` | Runs the application directly using `go run` |
| `make run-app` | Builds and launches `bin/OSXTrafficStats.app` using `open` |
| `make install` | Builds and installs `OSXTrafficStats.app` to `~/Applications` |
| `make test` | Runs unit tests |
| `make test-coverage` | Runs unit tests and produces an HTML coverage report |
| `make vet` | Runs `go vet` static analysis |
| `make clean` | Removes build and coverage artifacts |

## License

MIT