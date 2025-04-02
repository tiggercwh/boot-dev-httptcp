package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	IsClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		Listener: l,
	}
	s.IsClosed.Store(false)
	go func() {
		s.Listen()
	}()
	return s, nil
}

func (s *Server) Listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.IsClosed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" // Body
	conn.Write([]byte(response))
}

// Closes the listener and the server
func (s *Server) Close() error {
	s.IsClosed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}
