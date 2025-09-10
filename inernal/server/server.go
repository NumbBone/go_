package server

import (
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

type Handler func(w *responce.Writer, req *requests.Request) 

func runConnection(s *Server,conn io.ReadWriteCloser){
	defer conn.Close()
	
	respoceWiriter := responce.NewWirter(conn)

	r ,err := requests.ReqFromReader(conn)
	if err !=nil {
		respoceWiriter.WriteStatusLine(responce.BAD_REQUEST)
		respoceWiriter.WriteHeaders(*responce.GetDefaultHeaders(0))
		return
	}

	s.handler(respoceWiriter, r)
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

