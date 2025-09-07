package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Denis.test/inernal/requests"
	"Denis.test/inernal/responce"
	"Denis.test/inernal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *requests.Request) *server.HandlerError{
		if req.RequestLine.RequestTarget == "/yourproblem" {

			return &server.HandlerError{
				StatusCode: responce.BAD_REQUEST,
				Message: "Your problem is not my problem\n",
			}
		} else if req.RequestLine.RequestTarget == "/myproblem" {

			return &server.HandlerError{
				StatusCode: responce.INTERNAL_SERVER_ERROR,
				Message: "Woopsie, my bad\n",
			}
		} else {
			w.Write([]byte("All good, frfr\n"))
		}
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