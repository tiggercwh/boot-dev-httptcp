package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	res := Request{}
	b, err := io.ReadAll(reader)
	if err != nil {
		return &res, err
	}
	raw_splitted := strings.Split(string(b), "\r\n")
	splitted := strings.Split(string(raw_splitted[0]), " ")
	if len(splitted) != 3 {
		arg_err := fmt.Errorf("invalid amount of arguments")
		return &res, arg_err
	}
	method, req_target, http_v_str := splitted[0], splitted[1], splitted[2]
	http_v_splitted := strings.Split(http_v_str, "/")
	if method != strings.ToUpper(method) {
		meth_err := fmt.Errorf("method format error")
		return &res, meth_err
	}
	if len(http_v_splitted) <= 1 || http_v_splitted[1] != "1.1" {
		ver_err := fmt.Errorf("invalid http version")
		return &res, ver_err
	}
	res.RequestLine.HttpVersion = http_v_splitted[1]
	res.RequestLine.Method = method
	res.RequestLine.RequestTarget = req_target
	return &res, nil
}
