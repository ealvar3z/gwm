package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Action int

const (
	Close Action = iota
	SelectAbove
	SelectBelow
	ShrinkFront
	ExpandFront
)

type ActionKeyPress struct {
	Modifier uint16 `toml:"modifier"`
	Keysym   uint32 `toml:"keysym"`
	Action   Action `toml:"action"`
}

type Command struct {
	Modifier uint16 `toml:"modifier"`
	Keysym   uint32 `toml:"keysym"`
	Command  string `toml:"command"`
}

type Config struct {
	BorderThickness             uint32           `toml:"border_thickness"`
	BorderGap                   uint32           `toml:"border_gap"`
	ActiveBorder                uint32           `toml:"active_border"`
	InactiveBorder              uint32           `toml:"inactive_border"`
	WorkspaceModifier           uint16           `toml:"workspace_modifier"`
	WorkspaceMoveWindowModifier uint16           `toml:"workspace_move_window_modifier"`
	Autostart                   []string         `toml:"autostart"`
	Actions                     []ActionKeyPress `toml:"actions"`
	Commands                    []Command        `toml:"commands"`
}

func GetConfig() (*Config, error) {
	homeDir, exists := os.LookupEnv("HOME")
	if !exists {
		return &Config{}, fmt.Errorf("ERROR: no HOME var set.")
	}
	path := "%s/.config/gwm/config.toml"
	configDir := fmt.Sprintf(path, homeDir)

	configFile, err := os.Open(configDir)
	if err != nil {
		msg := "ERROR: Unable to open config.toml file from %s: %v"
		return &Config{}, fmt.Errorf(msg, configDir, err)
	}
	defer configFile.Close()

	var conf Config
	if _, err := toml.NewDecoder(configFile).Decode(&conf); err != nil {
		msg := "ERROR: Unable to decode toml config from %s: %v"
		return &Config{}, fmt.Errorf(msg, err)
	}

	return &conf, nil
}
