// +build win
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
		stubRelay := gpio.StubRelay{}
		stubRelay.SetConfig(relayConfig.ID, relayConfig.Name, relayConfig.Pin, relayConfig.Timings)
		stubRelay.SetOff()
		relays[idx] = &stubRelay
	}

	server.Init(relays)
}
