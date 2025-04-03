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
	n, err := w.Write(p)
	return n, err
}
