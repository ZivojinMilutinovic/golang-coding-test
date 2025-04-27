package main

import (
	"time"

	"github.com/ZivojinMilutinovic/golang-coding-test/client"
	"github.com/ZivojinMilutinovic/golang-coding-test/server"
)

func main() {
	go server.StartServer()
	// Give the server a moment to start before making requests
	time.Sleep(5 * time.Second) // Adjust if needed for your server to fully initialize
	client.TestClient()
	select {} //block forever
}
