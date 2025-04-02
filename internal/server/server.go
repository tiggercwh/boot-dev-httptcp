package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/bootdotdev/learn-http-protocol/internal/response"
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
	response.WriteStatusLine(conn, response.StatusOK)
	h := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, h)
}

// Closes the listener and the server
func (s *Server) Close() error {
	s.IsClosed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}
