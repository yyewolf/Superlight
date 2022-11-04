package main

import (
	"encoding/json"
	"os"
)

// Config : holds the config data
type Config struct {
	Host     string `json:"host"`
	ClientID string `json:"client_id"`
	Secret   string `json:"secret"`
	DeviceID string `json:"device_id"`
	Debug    bool   `json:"debug"`
}

var (
	config Config
)

// LoadConfig : loads superlight.json and returns any errors
func LoadConfig() error {
	data, err := os.ReadFile("superlight.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}

// InitialConfig : saves default superlight.json and returns any errors
func InitialConfig() error {
	config = Config{
		Host:     "https://openapi.tuyaeu.com",
		ClientID: "client_id",
		Secret:   "secret",
		DeviceID: "device_id",
		Debug:    false,
	}
	data, _ := json.MarshalIndent(config, "", "    ")
	err := os.WriteFile("superlight.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
