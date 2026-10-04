package main

import (
	"strings"
	"testing"
)

func TestGenerateLaunchAgentPlist(t *testing.T) {
	fakeExec := "/Applications/OSXTrafficStats.app/Contents/MacOS/osx-traffic-stats"
	plist := generateLaunchAgentPlist(fakeExec)

	if !strings.Contains(plist, "<string>com.smford.osx-traffic-stats</string>") {
		t.Errorf("plist missing expected label: %s", plist)
	}
	if !strings.Contains(plist, "<string>"+fakeExec+"</string>") {
		t.Errorf("plist missing executable path: %s", plist)
	}
	if !strings.Contains(plist, "<key>RunAtLoad</key>\n\t<true/>") {
		t.Errorf("plist missing RunAtLoad key: %s", plist)
	}
}

func TestLaunchAgentPath(t *testing.T) {
	path, err := launchAgentPath()
	if err != nil {
		t.Fatalf("launchAgentPath error: %v", err)
	}
	if !strings.HasSuffix(path, "Library/LaunchAgents/com.smford.osx-traffic-stats.plist") {
		t.Errorf("unexpected path: %s", path)
	}
}
