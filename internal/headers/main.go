package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}
	header := string(data[:idx])
	trimmed := strings.TrimSpace(header) //Host: localhost:42069, Host:localhost:42069
	splitted := strings.Split(trimmed, " ")
	first_split := splitted[0]
	colon_idx := strings.Index(first_split, ":")
	if len(splitted) > 2 {
		return 0, false, fmt.Errorf("invalid header")
	}
	if len(splitted) == 1 {
		if colon_idx != len(first_split)-1 && first_split[colon_idx+1] == ':' {
			return 0, false, fmt.Errorf("invalid header")
		}
		k := first_split[:colon_idx]
		v := first_split[colon_idx+1:]
		h[k] = v
	}
	if len(splitted) == 2 {
		if colon_idx != len(first_split)-1 {
			return 0, false, fmt.Errorf("invalid header")
		}
		k := first_split[:colon_idx]
		v := splitted[1]
		h[k] = v
	}
	fmt.Println(idx)
	return idx + 2, false, nil
}
