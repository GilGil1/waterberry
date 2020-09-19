// +build linux

package main

import (
	config "waterberry/internal/pkg/config"
	gpio "waterberry/internal/pkg/gpio"
	server "waterberry/internal/pkg/webserver"
)

var globalConfig config.GlobalConfig

func main() {
	config.LoadConfig(&globalConfig)
	gpio.Init(globalConfig)
	go gpio.StartIrrigation(globalConfig)
	server.Init()
}
