package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andresoro/hemera/pkg/backend"
	"github.com/andresoro/hemera/pkg/server"
)

// flag variables
var (
	graphitePort string
	serverHost   string
	serverPort   string
	purge        int
)

func init() {
	flag.StringVar(&graphitePort, "g", "2003", "graphite server port specifies where to purge metrics")
	flag.StringVar(&serverHost, "s", "localhost", "host for hemera server")
	flag.StringVar(&serverPort, "p", "8484", "port to listen for metrics")
	flag.IntVar(&purge, "t", 10, "interval for purging metrics to graphite in seconds")

	flag.Parse()
}

func main() {

	graphite := &backend.Graphite{Addr: fmt.Sprintf("localhost:%s", graphitePort)}

	purgeTime := time.Duration(purge) * time.Second

	srv, err := server.New(purgeTime, serverHost, serverPort, graphite)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening for metrics on port: %s \n", serverPort)
	log.Printf("purging metrics every %d seconds", purge)
	log.Printf("purging metrics to graphite server on port %s", graphitePort)

	srv.Run()

}
