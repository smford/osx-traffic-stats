package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type AppState struct {
	SessionSent       uint64          `json:"session_sent"`
	SessionRecv       uint64          `json:"session_recv"`
	TotalSent         uint64          `json:"total_sent"`
	TotalRecv         uint64          `json:"total_recv"`
	SelectedInterface string          `json:"selected_interface"`
	UnitMode          UnitMode        `json:"unit_mode"`
	StyleMode         StyleMode       `json:"style_mode"`
	MenuBarIconMode   MenuBarIconMode `json:"menubar_icon_mode,omitempty"`
	DataCapBytes      uint64          `json:"data_cap_bytes"`
}

var stateMu sync.Mutex

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "Application Support", "com.smford.osx-traffic-stats", "state.json"), nil
}

func loadAppState() (*AppState, error) {
	path, err := stateFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state AppState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func saveAppState(state *AppState) error {
	stateMu.Lock()
	defer stateMu.Unlock()

	path, err := stateFilePath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func persistCurrentState() {
	state := &AppState{
		SessionSent:       sessionSent.Load(),
		SessionRecv:       sessionRecv.Load(),
		TotalSent:         totalSent.Load(),
		TotalRecv:         totalRecv.Load(),
		SelectedInterface: getSelectedInterface(),
		UnitMode:          getUnitMode(),
		StyleMode:         getStyleMode(),
		MenuBarIconMode:   getMenuBarIconMode(),
		DataCapBytes:      getDataCap(),
	}
	_ = saveAppState(state)
}
