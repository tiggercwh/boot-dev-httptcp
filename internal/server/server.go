package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/bootdotdev/learn-http-protocol/internal/request"
	"github.com/bootdotdev/learn-http-protocol/internal/response"
)

type Server struct {
	Handler
	Listener net.Listener
	IsClosed atomic.Bool
}

type HandlerError struct {
	Message string
	response.StatusCode
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, herr *HandlerError) {
	if herr == nil {
		return
	}
	response.WriteStatusLine(w, herr.StatusCode)
	headers := response.GetDefaultHeaders(len(herr.Message))
	response.WriteHeaders(w, headers)
	w.Write([]byte(herr.Message))
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		Handler:  handler,
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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalln("fail to parse req")
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	herr := s.Handler(buf, req)
	if herr != nil {
		WriteHandlerError(conn, herr)
		return
	}
	response.WriteStatusLine(conn, response.StatusOK)
	h := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(conn, h)
	conn.Write(buf.Bytes())
}

// Closes the listener and the server
func (s *Server) Close() error {
	s.IsClosed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}
