package requests

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"Denis.test/inernal/headers"
)

type State string
const(
	initalized  State = "init"
	StateHeader State = "header"
	StateBody   State = "body"
	done        State = "done"
	errState    State = "error state"
)

func (r *Request) hasBody() bool  {
	// TODO: check for transfer-encoding: chunked
	length := getInt(r.Headers, "content-length", 0)
	return length > 0
}

func getInt(headers *headers.Headers, name string, defalVal int) int  {
	valueStr , exists := headers.Get(name)

	if !exists {
		return defalVal
	}

	val ,err := strconv.Atoi(valueStr)

	if err != nil {
		return defalVal
	}

	return val
}

type Request struct {
	RequestLine   RequestLine
	Headers       *headers.Headers
	Body 		  string
	ParserState   State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse (data []byte) (int , error){

	read := 0
outer:
for{
	curretData := data[read:]
	if len(curretData) == 0 {
	break outer
	}
	switch r.ParserState{
	case errState:
		return 0, Error_State
	case initalized:
		rl , n , err := parseRequestLine(curretData)
		if err != nil {
			r.ParserState = errState
			return 0, err
		}

		if n == 0{
			break outer
		}

		r.RequestLine = *rl
		read += n 
		
		r.ParserState = StateHeader
	case StateHeader:

		n , don , err := r.Headers.Parse(curretData)

		if err != nil {
			return 0, err
		}
		if n == 0 {
			break outer
		}

		if don {
			if r.hasBody() {
				r.ParserState = StateBody
			} else {
				r.ParserState = done
			}
		}
		read += n
		
	case StateBody:
		lengthstr := getInt(r.Headers, "content-length", 0)

		if lengthstr == 0 {
			panic("chunked not supported")
		}

		remain := min(lengthstr - len(r.Body) , len(curretData))
		r.Body += string(curretData[:remain])

		read += remain
	
		if len(r.Body) == lengthstr{
			r.ParserState = done
		}
	case done:
	break outer
	default:
		panic("How")
	}
}
  return read ,nil
}

func (r *Request) done() bool {
	return r.ParserState == done
}

func (r *Request) error() bool {
	return r.ParserState == errState
}

var ErrBadHeader = fmt.Errorf("bad start line")
var Error_State = fmt.Errorf("error sate in Request")
var LINESEP = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		ParserState: initalized,
		Headers: headers.NewHeaders(),
		Body: "",
	}
}

func parseRequestLine( s []byte) (*RequestLine,int, error)  {
	idx := bytes.Index(s,LINESEP);
	
	if idx == -1 {
		return nil, 0, nil;
	}

	startline := s[:idx]
	read := idx+len(LINESEP)

	parts := bytes.Split(startline, []byte(" "));

	if len(parts) != 3 {
		return nil ,0 ,ErrBadHeader
	}

	
	httparts := bytes.Split(parts[2], []byte("/"))
	
	if (len(httparts) != 2 || string(httparts[0]) != "HTTP" || string(httparts[1]) != "1.1") {
		return nil ,0, ErrBadHeader
	}
	
	rt := &RequestLine{
		Method: string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion: string(httparts[1]),
	}
	return rt, read ,nil 
}


func ReqFromReader(reader io.Reader) (*Request , error) {
	request := newRequest()

	buff := make([]byte,1024)
	buffidx := 0

	for !request.done() && !request.error(){
		n, err := reader.Read(buff[buffidx:])

		if err != nil{
			return nil ,err
		}
		
		buffidx += n
		ReadN , err := request.parse(buff[:buffidx])
		if err != nil {
			return nil ,err
		}

		copy(buff, buff[ReadN:buffidx])
		buffidx -= ReadN
	}

	return request ,nil
}