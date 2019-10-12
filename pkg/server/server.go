package server

import (
	"bytes"
	"log"
	"net"

	"github.com/andresoro/hemera/pkg/backend"
	"github.com/andresoro/hemera/pkg/cache"
	"github.com/andresoro/hemera/pkg/metric"
)

type Server struct {
	host    string
	port    string
	udpAddr *net.UDPAddr
	cache   *cache.Cache
	backend backend.Backend
}

func New(be backend.Backend, host, port string) (*Server, error) {

	service := host + ":" + port

	addr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		return nil, err
	}

	c := cache.New()

	server := &Server{
		host:    host,
		port:    port,
		udpAddr: addr,
		cache:   c,
		backend: be,
	}

	return server, nil
}

// Run UDP server on given addr
func (s *Server) Run() {

	ln, err := net.ListenUDP("udp", s.udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for packets on port %s \n", s.port)
	defer ln.Close()

	for {
		s.handleConn(ln)
	}

}

func (s *Server) handleConn(conn *net.UDPConn) {
	newLine := []byte("\n")
	buffer := make([]byte, 1024)

	// read data from conn onto our buffer up
	// n is an int of the number of bytes read
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("error %e", err)
	}

	payload := buffer[:n]

	// if we have more than one metric to record
	if bytes.Contains(payload, newLine) {
		payloads := bytes.Split(payload, newLine)

		for _, m := range payloads {
			newMetric, err := metric.Parse(m)
			if err != nil {
				continue
			}
			s.cache.Add(newMetric)
		}
	} else {
		m, err := metric.Parse(payload)
		// only push on successful parse
		if err == nil {
			s.cache.Add(m)
		}
	}

}
