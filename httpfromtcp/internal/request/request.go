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
	Method        string
	RequestTarget string
	HttpVersion   string
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

func validHTTP(HttpVersion, version string) bool {
	if strings.HasPrefix(HttpVersion, "HTTP/") && version == "1.1" {
		return true
	}
	return false
}

// Parsing of Request Line
func parseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, []byte(SEPERATOR))
	// NOTE: Return, If there is not SEPERATOR Found, /r/n
	if idx == -1 {
		return RequestLine{}, 0, nil
	}
	// NOTE: Found SEPERATOR on idx. So, Take a data before SEPERATOR
	line := string(data[:idx]) // Not Includes SEPERATOR
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line: expected 3 parts, got %d", len(parts))
	}

	method := parts[0]
	target := parts[1]
	versionRaw := parts[2]

	// NOTE: Validate method: must be all uppercase letters
	if !methodRegex.MatchString(method) {
		return RequestLine{}, 0, fmt.Errorf("invalid method: %q", method)
	}

	const prefix = "HTTP/"

	// Gives version of HTTP. Like 1.1 or 2, whatever
	version := strings.TrimPrefix(versionRaw, prefix)

	ok := validHTTP(versionRaw, version)
	if ok != true {
		return RequestLine{}, 0, fmt.Errorf("unsupported http version or format: %q", version)
	}

	rl := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}
	return rl, idx + len(SEPERATOR), nil
}

const bufSize = 30

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := newRequest()
	buf := make([]byte, bufSize)
	bufLen := 0

	// keep reading until request is fully parsed
	for !r.isDone() {
		// NOTE: It reads from starting, So At start, bufLen is 0. So it will reads all thing from starting which is 0
		n, err := reader.Read(buf[bufLen:]) // NOTE: Read all thing from bufLen
		if err != nil {
			return nil, err
		}

		bufLen += n
		// NOTE: Parse until bufLen, which is n at first
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

// NOTE: Check If body available for not, using Content-Length
func (r *Request) hasBody() bool {
	contentLength := getInt(r.Headers, "Content-Length", 0)
	return contentLength > 0
}

// NOTE: This is a main function, Where Everything get attached. Like Statemachine
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
			// NOTE: 1st Phase, Parse request line
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
			// NOTE Going for 2nd stage, which is parsing headers
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
