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
	server.BuildTime = BuildTime
	server.GitCommit = GitCommit
	config := config.LoadConfig("config.json")
	relays := make([]gpio.IRelay, 4)

	// Set Relay Config
	for idx, relayConfig := range config.Relays {
		var relay gpio.IRelay = gpio.NewGPIORelay()
		relays[idx] = relay
		relay.SetConfig(relayConfig.ID, relayConfig.Name, relayConfig.Pin, relayConfig.Timings)
		relay.Init()
	}

	server.Init(relays, config)

}
