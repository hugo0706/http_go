package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type status int

const (
	initialized status = iota
	parsingHeaders
	done
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{state: initialized}
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	for request.state != done {
		if readToIndex >= len(buffer){
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}
		n, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = done
				break
			}
			return nil, err
		}
		readToIndex += n
		
		n, err = request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}
		if n > 0 {
			newBuf := make([]byte, len(buffer)-n)
			copy(newBuf, buffer[n:])
			buffer = newBuf
			readToIndex -= n
		}
	}
	
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	n := 0
	switch r.state {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil || n == 0 {
			return 0, err
		}
		r.state = parsingHeaders
		r.RequestLine = *requestLine
	case parsingHeaders:
		//TODO
	case done:
		return 0, errors.New("trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
	
	return n, nil
}

func parseRequestLine(content []byte) (*RequestLine, int, error) {
	parts := strings.Split(string(content), crlf)
	if len(parts) == 1 {
		return nil, 0, nil
	}

	requestLineData := parts[0]
	requestLine, err := requestLineFromString(requestLineData)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, len(requestLineData), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	requestLineParts := strings.Split(str, " ")
	requestLine := &RequestLine{}

	if len(requestLineParts) != 3 {
		return nil, errors.New("Invalid request line")
	}
	requestLine.Method = requestLineParts[0]
	if ok, _ := regexp.Match("^[A-Z]+$", []byte(requestLine.Method)) ; !ok {
		return nil, errors.New("Invalid method format")
	}
	if requestLineParts[2] != "HTTP/1.1" {
		return nil, errors.New("Only HTTP/1.1 supported")
	}
	requestLine.HttpVersion = strings.Split(requestLineParts[2], "/")[1]
	requestLine.RequestTarget = requestLineParts[1]
	
	return requestLine, nil
}