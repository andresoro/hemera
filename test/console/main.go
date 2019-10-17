package main

import (
	"fmt"
	"log"
	"time"

	"github.com/andresoro/hemera/pkg/cache"
	"github.com/andresoro/hemera/pkg/server"
)

type ConsoleBackend struct{}

func main() {

	backend := ConsoleBackend{}

	srv, err := server.New(5*time.Second, "localhost", "8080", backend)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run()

}

// Purge cache
func (cb ConsoleBackend) Purge(c *cache.Cache) error {

	// only handle counters for this example
	for name, value := range c.Counters {
		fmt.Printf("%s %f \n", name, value)
	}

	return nil
}
