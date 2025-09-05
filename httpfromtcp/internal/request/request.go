package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var SEPERATOR = "\r\n"

var (
	methodRegex = regexp.MustCompile(`^[A-Z]+$`)
)

type ParserState string

const (
	InitState ParserState = "initialized"
	DoneState ParserState = "done"
)

func newRequest() *Request {
	request := &Request{
		State: InitState,
	}
	return request
}

// Parsing of Request Line
func parseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, []byte(SEPERATOR))
	if idx == -1 {
		return RequestLine{}, 0, nil
	}
	line := string(data[:idx])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line: expected 3 parts, got %d", len(parts))
	}

	method := parts[0]
	target := parts[1]
	versionRaw := parts[2]

	// Validate method: must be all uppercase letters
	if !methodRegex.MatchString(method) {
		return RequestLine{}, 0, fmt.Errorf("invalid method: %q", method)
	}

	// Validate HTTP version
	const prefix = "HTTP/"
	if !strings.HasPrefix(versionRaw, prefix) {
		return RequestLine{}, 0, fmt.Errorf("invalid http version format: %q", versionRaw)
	}
	version := strings.TrimPrefix(versionRaw, prefix)
	if version != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("unsupported http version: %q", version)
	}

	rl := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}
	return rl, idx + len(SEPERATOR), nil
}

const bufSize = 1024

// Get Request from memory using io.ReadAll
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := newRequest()
	buf := make([]byte, bufSize)
	bufLen := 0
	for !r.isDone() {
		n, err := reader.Read(buf[bufLen:]) // NOTE: Read all thing from bufLen
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := r.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return r, nil
}

func (r *Request) isDone() bool {
	return r.State == DoneState
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case InitState:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = rl
			read += n
			r.State = DoneState
		case DoneState:
			break outer
		}
	}

	return read, nil
}
