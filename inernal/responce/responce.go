package responce

import (
	"fmt"
	"io"

	"Denis.test/inernal/headers"
)

type StatusCode int

const (
	OK                    StatusCode = 200
	BAD_REQUEST           StatusCode = 400
	INTERNAL_SERVER_ERROR StatusCode = 500
)

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}


type Writer struct{
	writer io.Writer
}

func NewWirter(writer io.Writer) *Writer{
	return &Writer{writer: writer}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error{
	
	statLine := []byte{}
	switch statusCode{
	case OK: statLine = []byte("HTTP/1.1 200 OK\r\n")
	case BAD_REQUEST: statLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case INTERNAL_SERVER_ERROR: statLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	}

	_, err := w.writer.Write(statLine)

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
		var err error = nil
	b := []byte{}
	headers.ForEach(func(name, value string) {
		b = fmt.Appendf(b, "%s: %s\r\n", name, value)
	})

	b = fmt.Append(b, "\r\n")

	_,err = w.writer.Write(b)

	return err
}
func (w *Writer) WriteBody(p []byte) (int, error){
	
	n ,err := w.writer.Write(p)

	return  n , err
}