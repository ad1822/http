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

type ParserState string

const (
	InitState ParserState = "initialized"
	DoneState ParserState = "done"
)

var SEPERATOR = "\r\n"

var (
	methodRegex = regexp.MustCompile(`^[A-Z]+$`)
)

// Get Request from memory using io.ReadAll
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{State: InitState}
	buf := make([]byte, 0, 8) // start with small buffer
	tmp := make([]byte, 8)    // read chunks into this
	for {
		n, err := reader.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			for {
				consumed, perr := r.parse(buf)
				if perr != nil {
					return nil, perr
				}
				if consumed == 0 {
					break // need more data
				}
				buf = buf[consumed:] // discard consumed bytes
				if r.State == DoneState {
					return r, nil
				}
			}
		}
		if err == io.EOF {
			if r.State == DoneState {
				return r, nil
			}
			return nil, fmt.Errorf("incomplete request")
		}
		if err != nil {
			return nil, err
		}
	}
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

func (r *Request) parse(data []byte) (int, error) {
	if r.State == DoneState {
		return 0, nil
	}

	rl, consumed, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if consumed == 0 {
		// not enough data yet
		return 0, nil
	}

	r.RequestLine = rl
	r.State = DoneState
	return consumed, nil
}
