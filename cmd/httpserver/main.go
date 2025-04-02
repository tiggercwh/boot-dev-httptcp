package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-http-protocol/internal/request"
	"github.com/bootdotdev/learn-http-protocol/internal/response"
	"github.com/bootdotdev/learn-http-protocol/internal/server"
	s "github.com/bootdotdev/learn-http-protocol/internal/server"
)

const port = 42069

func main() {
	server, err := s.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{
				Message:    "Your problem is not my problem\n",
				StatusCode: response.StatusBadRequest,
			}
		}
		if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{
				Message:    "Woopsie, my bad\n",
				StatusCode: response.StatusInternalError,
			}
		}
		w.Write([]byte("All good, frfr\n"))
		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
