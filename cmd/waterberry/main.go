// +build linux
package main

import (
	config "waterberry/internal/pkg/config"
	gpio "waterberry/internal/pkg/gpio"
	server "waterberry/internal/pkg/web"
)

var GitCommit = "Unknown"
var BuildTime = "Unknown"

func main() {
	config := config.LoadConfig("config.json")
	relays := make([]gpio.IRelay, 4)

	// Set Relay Config
	for idx, relayConfig := range config.Relays {
		gpioRelay := gpio.GPIORelay{}
		gpioRelay.SetConfig(relayConfig.ID, relayConfig.Name, relayConfig.Pin, relayConfig.Timings)
		relays[idx] = &gpioRelay
	}

	server.Init(relays)
	server.BuildTime = BuildTime
	server.GitCommit = GitCommit
}
