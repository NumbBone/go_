package server

import (
	"fmt"
	"io"
	"net"

	"Denis.test/inernal/responce"
)

type Server struct {
	closed bool
}

func runConnection(_s *Server,conn io.ReadWriteCloser){
	defer conn.Close()
	
	buf := make([]byte, 4096) 
    _, err := conn.Read(buf)
    if err != nil && err != io.EOF {
        return
    }

	header := responce.GetDefaultHeaders(0)
	responce.WriteStatusLine(conn , responce.OK)
	responce.WriteHeaders(conn , header)
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

func Serve(port int) (*Server, error) {
	lisener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil ,err 
	}

	serve := &Server{ closed: false }
	go runServer(serve, lisener)

	return serve, nil
}

func (s *Server) Close() error{
	s.closed = true
	return nil
}
