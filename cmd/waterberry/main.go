// +build linux
package main

import (
	config "waterberry/internal/pkg/config"
	gpio "waterberry/internal/pkg/gpio"
	server "waterberry/internal/pkg/web"
)

func main() {
	config := config.LoadConfig("config.json")
	relays := make([]gpio.IRelay, 4)

	// Set Relay Config
	for idx, relayConfig := range config.Relays {
		GPIORelay := gpio.GPIORelay{
			ID:            relayConfig.ID,
			Pin:           relayConfig.Pin,
			Name:          relayConfig.Name,
			CurrentStatus: gpio.RelayOff,
		}
		relays[idx] = &GPIORelay
	}

	server.Init(relays)
}
