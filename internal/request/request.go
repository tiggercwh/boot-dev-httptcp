package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ParserState int

const (
	Initialized ParserState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		ParserState: Initialized,
	}
	var err error
	for req.ParserState != Done {
		if readToIndex >= bufferSize {
			new_buf := make([]byte, len(buf)*2)
			copy(new_buf, buf)
			buf = new_buf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.ParserState = Done
			break
		}
		readToIndex += n
		_, err = req.parse(buf)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("req", req)
	return req, err
}

func parseRequestLine(data []byte) (int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil
	}
	return idx, nil
}

// Implement a new func (r *Request) parse(data []byte) (int, error) method.
// It accepts the next slice of bytes that needs to be parsed into the Request struct
// It updates the "state" of the parser, and the parsed RequestLine field.
// It returns the number of bytes it consumed (meaning successfully parsed) and an error if it encountered one.
func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == Initialized {
		bytes_len, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytes_len == 0 {
			return 0, nil
		}
		requestLine, err := requestLineFromString(string(data[:bytes_len]))
		if err != nil {
			return 0, err
		}
		r.RequestLine = *requestLine
		r.ParserState = Done
		return bytes_len, nil
	}
	return 0, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
