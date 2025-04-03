package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

const crlf = "\r\n"
const specials = "!#$%&'*+-.^_`|~"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	low_key := strings.ToLower(key)

	badchar_key := strings.ContainsFunc(key, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune(specials, r))
	})

	if badchar_key {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	h.Set(low_key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}
