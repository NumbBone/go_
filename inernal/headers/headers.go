package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var endOfLine = []byte("\r\n")

func isToken(str []byte) bool{
	for _, b := range str {
		if b <= 32 || b >= 127 || strings.ContainsRune("()<>@,;:\\\"/[]?={} \t", rune(b)) {
			return false
		}
	}
	return true
}

func parseheader(fieldline []byte) (string, string, error) {
	
	pr := bytes.TrimSpace(fieldline)

	parts:= bytes.SplitN(pr, []byte(":"), 2)
	if len(parts) != 2{
		return "","",fmt.Errorf("malformed header ")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name ,[]byte(" ")){
		return "","",fmt.Errorf("malformed header name")
	}

	return string(name) , string(value), nil
}


type Headers struct {
	headers  map[string]string}


func NewHeaders() *Headers  {
	return &Headers{
		headers :map[string]string{}}
}

func (h *Headers) Get(name string) (string ,bool) {
	str , ok := h.headers[strings.ToLower(name)]

	return str ,ok 
}

func (h *Headers) Delete(name string){
	name = strings.ToLower(name)
	delete(h.headers , name)
}

func (h *Headers) Set(name ,value string) {
	name = strings.ToLower(name)
	
	if v ,s := h.headers[name]; s {
		h.headers[name] = fmt.Sprintf("%s,%s", v , value)
	} else {
		h.headers[name] = value
	}

	
}


func (h *Headers) Replace(name ,value string) {
	name = strings.ToLower(name)
		h.headers[name] = value
}

func (h *Headers) ForEach(f func(name, value string)) {
	for k, v := range h.headers {
		f(k, v)
	}
}


func (h *Headers) Parse(data []byte) ( int, bool, error) {

	read := 0
	done := false

	for{
		idx := bytes.Index(data[read:], endOfLine)

		if idx == -1 {
			break 
		}

		if idx == 0 {
			done = true
			read += len(endOfLine)
			break
		}

		name ,value ,err:= parseheader(data[read:read+idx])

		if err != nil {
			return 0, false ,err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}

		h.Set(name, value) 

		read += idx + len(endOfLine)
	}

	return read, done, nil
}