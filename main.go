package main

import (
	"log"
)

func main() {
	port, err := returnFirstAvailablePort(3000, 20)
	if err != nil {
		log.Fatalln(err)
	}
	server := NewServer(port)
	log.Fatal(server.Start())
}
