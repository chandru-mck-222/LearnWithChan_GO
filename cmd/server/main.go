package main

import (
	"chanlearning/pkg/service/server"
	"log"
)

func main() {
	if err := server.SetupServer(); err != nil {
		log.Fatal(err)
	}
}
