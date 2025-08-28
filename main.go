package main

import (
	"fmt"
	"net"

	"Denis.test/inernal/requests"
)


func main() {
	lisener ,err := net.Listen("tcp", ":42069");
	fmt.Println("Listening on port 42069")
	if err != nil {
		fmt.Println("Error listen:", err)
		return
	}

	
	for {
		conn, err := lisener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		r ,err := requests.ReqFromReader(conn)
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		fmt.Print("Request line\n")
		fmt.Print("- Method: ",r.RequestLine.Method,"\n")
		fmt.Print("- Target: ",r.RequestLine.RequestTarget,"\n")
		fmt.Print("- Version: ",r.RequestLine.HttpVersion,"\n")
		fmt.Print("Headers\n")
		r.Headers.ForEach(func(name, value string) {
			fmt.Printf("- %s: %s\n", name, value)
		})
		fmt.Print("- Body: ", string(r.Body),"\n")
	}
}