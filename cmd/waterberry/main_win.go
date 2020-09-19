// +build win
package main

import (
	config "waterberry/internal/pkg/config"
	server "waterberry/internal/pkg/webserver"
)

var globalConfig config.GlobalConfig

func main() {
	config.LoadConfig(&globalConfig)
	server.Init()
}
