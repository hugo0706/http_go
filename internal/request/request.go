package request

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hugo0706/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
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
	req := &Request{
		state: initialized,
		Headers: headers.Headers{},
	}
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	for req.state != done {
		if readToIndex >= len(buffer){
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != done {
					return nil, fmt.Errorf("incomplete request, in state: %d, read %d bytes on EOF", req.state, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		
		numBytesParsed, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed	
	}
	
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
		fmt.Printf("Bytes parsed: %d\n", totalBytesParsed)
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, err
		}
		r.state = parsingHeaders
		r.RequestLine = *requestLine
		return n, nil
	case parsingHeaders:
		n, finish, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if finish {
			r.state = done
		}
		return n, nil
	case done:
		return 0, errors.New("trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
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

	return requestLine, len(requestLineData) + 2, nil
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