package main

import "log"

func main() {
	server := newServer()
	server.setUpRoutes()
	log.Fatal(server.run())
}
