package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var config Config

/*
Config represents the structure of the configuration file.
*/
type Config struct {
	BaseURLPath string              `json:"base_url_path"`
	BindAddress string              `json:"bind_address"`
	// Targets maps target names to their respective URLs
	Targets     map[string]struct {
		URL string `json:"url"`
	}   `json:"targets"`
	// Groupings define target groupings for different types of activities
	Groupings   map[string][]string `json:"groupings"`
	DefaultGroup string             `json:"default_group"`
}


func loadConfig(location string) error {
	fmt.Printf("Loading configuration from %s\n", location)
	// Read the file
	data, err := os.ReadFile(location)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	// Parse the JSON
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.BaseURLPath == "" {
		return fmt.Errorf("base_url_path is required in the config file")
	}
	if len(config.Targets) == 0 {
		return fmt.Errorf("at least one target is required in the config file")
	}
	if config.DefaultGroup == "" {
		return fmt.Errorf("default_group is required in the config file")
	}
	if len(config.Groupings) == 0 {
		return fmt.Errorf("at least one grouping is required in the config file")
	}

	// Set default values if necessary
	if config.BindAddress == "" {
		config.BindAddress = ":8080"
		fmt.Println("bind_address not set in config, using default :8080")
	}

	fmt.Println("Configuration loaded successfully")

	return nil
}