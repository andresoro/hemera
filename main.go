package main

import (
	"log"
	"time"

	"github.com/andresoro/hemera/pkg/backend"
	"github.com/andresoro/hemera/pkg/server"
)

func main() {

	graphite := &backend.Graphite{Addr: "localhost:2023"}

	srv, err := server.New(10*time.Second, "localhost", "8484", graphite)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run()

}
