package response

import (
	"fmt"
	"io"

	"github.com/ad1822/httpfromtcp/internal/headers"
)

type Response struct {
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) *headers.Headers {
	headers := headers.NewHeaders()
	// NOTE: Content-Length will set up from message
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text")

	return headers
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOk:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")

	default:
		return fmt.Errorf("Unrecognized Format")
	}

	_, err := w.writer.Write(statusLine)
	return err
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	return n, err
}
