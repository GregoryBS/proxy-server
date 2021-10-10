package main

import (
	"proxy/pkg/api"
	"proxy/pkg/server"
)

func main() {
	go api.Run()
	server.Run(":8080")
}
