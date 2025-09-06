package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/ad1822/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	State       ParserState
	Body        string
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
	StateInit    ParserState = "initialized"
	StateDone    ParserState = "done"
	StateHeaders ParserState = "headers"
	StateBody    ParserState = "body"
)

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := headers.Get(name)

	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request {
	request := &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
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
	return r.State == StateDone
}

func (r *Request) hasBody() bool {
	contentLength := getInt(r.Headers, "Content-Length", 0)
	return contentLength > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.State {
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = rl
			read += n
			// NOTE: Instead of Done, Use this state to parse headers
			r.State = StateHeaders

		// NOTE: After start-line, It'll parse headers (field-line)
		case StateHeaders:
			n, done, err := r.Headers.Parse([]byte(currentData))
			if err != nil {
				return 0, err
			}

			// NOTE: Got a Full header
			if n == 0 {
				break outer
			}

			read += n
			if done {
				if r.hasBody() {
					r.State = StateBody
				} else {
					r.State = StateDone
				}
			}

		// NOTE: After headers, It'll parse Body
		case StateBody:
			contentLength := getInt(r.Headers, "Content-Length", 0)
			if contentLength == 0 {
				r.State = StateDone
				break
			}

			remaining := min(contentLength-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			// fmt.Printf("%s\n", r.Body)

			if contentLength == len(r.Body) {
				r.State = StateDone
			}

		// NOTE: Done. Everything parsed
		case StateDone:
			break outer
		default:
			panic("POOR Program")
		}
	}

	return read, nil
}
