package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Denis.test/inernal/headers"
	"Denis.test/inernal/requests"
	"Denis.test/inernal/responce"
	"Denis.test/inernal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w responce.Writer, req *requests.Request) {
		headers := responce.GetDefaultHeaders(0)

		if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(responce.BAD_REQUEST)
			w.WriteHeaders(*headers)
		} else if req.RequestLine.RequestTarget == "/myproblem" {

			}
		} else {
			
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