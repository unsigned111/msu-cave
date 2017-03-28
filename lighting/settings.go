package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Settings is a structure that holds settings values.
type Settings struct {
	ListenIP, ControlIP        string
	ListenPort, ControlPort    int

	StartColor, EndColor       LightingColor
	OffStartColor, OffEndColor LightingColor
	PulseAmount                float32
	PulseLength, PulsePause    int
	DefaultOn                  bool
	HeadsetChannel             int
	Scaler                     int
	FPS                        int
}

// ParseConfig parses the settings values from a JSON file path.
func ParseConfig(path string) (settings Settings) {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read the config file: %s", path)
	}
	err = json.Unmarshal(config, &settings)
	return settings
}
