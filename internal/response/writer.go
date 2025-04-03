package response

import (
	"fmt"
	"io"

	"github.com/bootdotdev/learn-http-protocol/internal/headers"
)

type WriterState int

const (
	WritingStatusLine WriterState = iota
	WritingHeaders
	WritingBody
	Done
)

type Writer struct {
	io.Writer
	WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer:      w,
		WriterState: WritingStatusLine,
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WritingHeaders {
		return fmt.Errorf("invalid writer state: %v", w.WriterState)
	}

	defer func() {
		w.WriterState = WritingBody
	}()
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != WritingStatusLine {
		return fmt.Errorf("invalid writer state: %v", w.WriterState)
	}

	defer func() {
		w.WriterState = WritingHeaders
	}()
	_, err := w.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState != WritingBody {
		return 0, fmt.Errorf("invalid writer state: %v", w.WriterState)
	}
	defer func() {
		w.WriterState = Done
	}()
	n, err := w.Write(p)
	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	var total int

	// 1. Write chunk size in hex + \r\n
	header := fmt.Sprintf("%X\r\n", len(p))
	n, err := w.Write([]byte(header))
	total += n
	if err != nil {
		return total, err
	}

	// 2. Write the actual data
	n, err = w.Write(p)
	total += n
	if err != nil {
		return total, err
	}

	// 3. Write trailing \r\n
	n, err = w.Write([]byte("\r\n"))
	total += n
	return total, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	// Write the final 0-length chunk: "0\r\n\r\n"
	done := []byte("0\r\n\r\n")
	return w.Write(done)
}

// func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
// 	if w.WriterState != WritingBody {
// 		return 0, fmt.Errorf("invalid writer state: %v", w.WriterState)
// 	}

// 	chunk_len := len(p)
// 	if len(p) == 0 {
// 		return w.WriteChunkedBodyDone()
// 	}
// 	hex_len := fmt.Sprintf("%X", chunk_len)
// 	hexlen_bytes := []byte(hex_len + "\r\n")
// 	n1, err := w.Write(hexlen_bytes)
// 	if err != nil {
// 		return 0, err
// 	}
// 	n2, err := w.Write(p)
// 	if err != nil {
// 		return n1, err
// 	}
// 	return n1 + n2, err
// }

// func (w *Writer) WriteChunkedBodyDone() (int, error) {
// 	defer func() {
// 		w.WriterState = Done
// 	}()
// 	n, err := w.Write([]byte("0\r\n\r\n"))
// 	return n, err
// }
