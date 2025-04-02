package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	s "github.com/bootdotdev/learn-http-protocol/internal/server"
)

const port = 42069

func main() {
	server, err := s.Serve(port)
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
