// +build win
package main

import (
	"net/http"
	config "waterberry/internal/pkg/config"
	gpio "waterberry/internal/pkg/gpio"
	server "waterberry/internal/pkg/web"

	"github.com/markbates/pkger"
)

func main() {
	config := config.LoadConfig("config.json")
	relays := make([]gpio.IRelay, 4)

	// Set Relay Config
	for idx, relayConfig := range config.Relays {
		stubRelay := gpio.StubRelay{
			ID:            relayConfig.ID,
			Pin:           relayConfig.Pin,
			Name:          relayConfig.Name,
			CurrentStatus: gpio.RelayOff,
		}
		relays[idx] = &stubRelay
	}
	dir := http.FileServer(pkger.Dir("/public"))
	http.ListenAndServe(":3000", dir)

	server.Init(relays)
}
