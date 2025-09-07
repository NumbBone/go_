package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"Denis.test/inernal/requests"
	"Denis.test/inernal/responce"
)

type Server struct {
	closed bool
	handler Handler
}

type HandlerError struct{
	StatusCode  responce.StatusCode
	Message 	string
}

func runConnection(s *Server,conn io.ReadWriteCloser){
	defer conn.Close()
	

	header := responce.GetDefaultHeaders(0)
	r ,err := requests.ReqFromReader(conn)
	if err !=nil {
		responce.WriteStatusLine(conn , responce.BAD_REQUEST)
		responce.WriteHeaders(conn , header)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	errhandle := s.handler(writer, r)

	var body []byte = nil
	var sat responce.StatusCode = responce.OK

	if errhandle != nil {
		sat = errhandle.StatusCode
		body = []byte(errhandle.Message)
	} else {
		body = writer.Bytes()
			
	}

	header.Replace("Content-Length", fmt.Sprintf("%d",len(body)))

	responce.WriteStatusLine(conn , sat)
	responce.WriteHeaders(conn , header)
	conn.Write(body)
}

func runServer(s *Server, listener net.Listener)  {
	for{
		conn , err := listener.Accept()

		if s.closed {
			return
		}

		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port int , handler Handler) (*Server, error) {
	lisener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil ,err 
	}

	serve := &Server{ 
		closed: false,
		handler: handler,
	}

	go runServer(serve, lisener)

	return serve, nil
}

func (s *Server) Close() error{
	s.closed = true
	return nil
}

type Handler func(w io.Writer, req *requests.Request) *HandlerError