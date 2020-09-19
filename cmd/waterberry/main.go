// +build linux

package main

import (
	gpio "waterberry/internal/pkg/gpio"
	server "waterberry/internal/pkg/webserver"
)

func main() {

	go gpio.Init()
	server.Init()

}
