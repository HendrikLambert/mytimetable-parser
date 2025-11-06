package main

import (
	"fmt"
	"os"
)

func main() {
	confLoc := os.Getenv("CONFIG_LOCATION")

	if confLoc == "" {
		confLoc = "config.json"
		fmt.Println("CONFIG_LOCATION environment variable not provided, using default config.json")
	}

	err := loadConfig(confLoc)
	
	if err != nil {
		fmt.Printf("Error loading config from %s: %v\n", confLoc, err)
		return
	}

}