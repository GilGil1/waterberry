package gpio

// +build linux

import (
	"fmt"
	"os"
	"strings"
	"time"

	config "waterberry/internal/pkg/config"

	"github.com/stianeikeland/go-rpio"
)

var globalConfig config.GlobalConfig

func Init(config config.GlobalConfig) {

	globalConfig = config
	setupRalays(config)
}

func StartIrrigation(config config.GlobalConfig) {
	for {
		time.Sleep(time.Second)
	}
}

func setupRalays(config config.GlobalConfig) {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	for i, relayConfig := range config.Relays {
		fmt.Printf("relay: %d, pin=%d, config =%v\n", i, relayConfig.Pin, relayConfig)
		pin := rpio.Pin(relayConfig.Pin)
		pin.Output()
		pin.Low()
		fmt.Println("Pins are set")
	}
}

func SetRelayMode(relay int, mode string) {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	fmt.Printf("relay = %d, mode = %s\n", relay, mode)
	relayConfig := globalConfig.Relays[relay]
	fmt.Printf("%v\n", relayConfig)
	pin := rpio.Pin(relayConfig.Pin)
	pin.Output()
	switch strings.ToLower(mode) {
	case "on":
		pin.Low()

	case "off":
		pin.High()

	case "toggle":
		pin.Toggle()

	default:
		fmt.Println("mode %s undupported\n", mode)
	}
}
