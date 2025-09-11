package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Denis.test/inernal/requests"
	"Denis.test/inernal/responce"
	"Denis.test/inernal/server"
)

const port = 42069

func responce400() []byte{
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func responce500() []byte{
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func responce200() []byte{
	return []byte(`<html>
  <head>
	<title>200 OK</title>
  </head>
  <body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	server, err := server.Serve(port, func(w *responce.Writer, req *requests.Request) {
		headers := responce.GetDefaultHeaders(0)
		body := responce200()
		stat := responce.OK

		
		if req.RequestLine.RequestTarget == "/yourproblem" {

			stat = responce.BAD_REQUEST
			body = responce400()

		} else if req.RequestLine.RequestTarget == "/myproblem" {

			stat = responce.INTERNAL_SERVER_ERROR
			body = responce500()
			
		}

			w.WriteStatusLine(stat)
			headers.Replace("content-length",fmt.Sprintf("%d",len(body)))
			headers.Replace("Content-Type", "text/plain")
			w.WriteHeaders(*headers)
			w.WriteBody(body)

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