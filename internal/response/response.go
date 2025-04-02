package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/bootdotdev/learn-http-protocol/internal/headers"
)

type StatusCode int

const (
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	base_str := fmt.Sprintf("HTTP/1.1 %v ", statusCode)
	var err error
	switch statusCode {
	case StatusOK:
		base_str += "OK\r\n"
	case StatusBadRequest:
		base_str += "Bad Request\r\n"
	case StatusInternalError:
		base_str += "Internal Server Error\r\n"
	default:
		base_str += "\r\n"
	}
	_, err = w.Write([]byte(base_str))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	clen := strconv.Itoa(contentLen)
	h.Set("Content-Length", clen)
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		res := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(res))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
