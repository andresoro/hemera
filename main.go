package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/andresoro/hemera/pkg/backend"
	"github.com/andresoro/hemera/pkg/server"
)

type config struct {
	graphitePort string
	serverHost   string
	serverPort   string
	purge        time.Duration
}

func configure() *config {
	var cfg config
	flag.StringVar(&cfg.graphitePort, "g", "2003", "graphite server port specifies where to purge metrics")
	flag.StringVar(&cfg.serverHost, "s", "localhost", "host for hemera server")
	flag.StringVar(&cfg.serverPort, "p", "8484", "port to listen for metrics")
	flag.DurationVar(&cfg.purge, "t", 10*time.Second, "interval for purging metrics to graphite in seconds")

	flag.Parse()
	return &cfg
}

func main() {
	cfg := configure()
	graphite := &backend.Graphite{Addr: fmt.Sprintf("localhost:%s", cfg.graphitePort)}

	srv, err := server.New(cfg.purge, cfg.serverHost, cfg.serverPort, graphite)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening for metrics on port: %s \n", cfg.serverPort)
	log.Printf("purging metrics every %d seconds", cfg.purge)
	log.Printf("purging metrics to graphite server on port %s", cfg.graphitePort)

	srv.Run()
}
