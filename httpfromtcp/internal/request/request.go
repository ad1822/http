package request

import (
	"fmt"
	"io"
	"regexp"
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

var (
	methodRegex = regexp.MustCompile(`^[A-Z]+$`)
)

// Get Request from memory using io.ReadAll
func RequestFromReader(r io.Reader) (*Request, error) {
	line, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("Error in Read from Reader : ", err)
	}

	s := string(line)
	lines := strings.SplitN(s, "\r\n", 2)

	parsedLine, err := parseRequestLine(lines[0])
	if err != nil {
		fmt.Println("Error in parsing : ", err)
	}

	return &Request{RequestLine: parsedLine}, nil
}

// Parsing of Request Line
func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line: expected 3 parts, got %d", len(parts))
	}

	method := parts[0]
	target := parts[1]
	versionRaw := parts[2]

	// Validate method: must be all uppercase letters
	if !methodRegex.MatchString(method) {
		return RequestLine{}, fmt.Errorf("invalid method: %q", method)
	}

	// Validate HTTP version
	const prefix = "HTTP/"
	if !strings.HasPrefix(versionRaw, prefix) {
		return RequestLine{}, fmt.Errorf("invalid http version format: %q", versionRaw)
	}
	version := strings.TrimPrefix(versionRaw, prefix)
	if version != "1.1" {
		return RequestLine{}, fmt.Errorf("unsupported http version: %q", version)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}
