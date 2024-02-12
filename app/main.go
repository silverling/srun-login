package main

import (
	"log"
	"path/filepath"
	"time"
)

func run() {
	configPath := filepath.Join(GetProgramFolder(), "config.yaml")
	config, err := loadConfig(configPath)
	if err != nil {
		log.Output(2, err.Error())
		return
	}

	isOnline := testConnection()
	if isOnline {
		log.Println("You are already online")
	}

	for {
		isOnline = testConnection()
		if !isOnline {
			log.Println("Offline, try to login")
			for {
				isLoggedin := login(config)
				if isLoggedin {
					break
				} else {
					time.Sleep(5 * time.Second)
				}
			}
		} else {
			time.Sleep(60 * time.Second)
		}
	}
}

func main() {
	RunService(run)
}
