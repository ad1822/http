package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var SEPERATOR = []byte("\r\n")

func NewHeaders() Headers {

	return map[string]string{}
}

func parseHeaders(fieldLine []byte) (string, string, error) {

	parts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}
	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	done := false
	read := 0
	for {
		idx := bytes.Index(data[read:], SEPERATOR)
		// NOTE: Seperator not found
		if idx == -1 {
			break
		}

		// NOTE: We got a full header
		if idx == 0 {
			read += len(SEPERATOR)
			done = true
			break
		}

		key, value, err := parseHeaders(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		read += idx + len(SEPERATOR)
		h[key] = value
	}
	return read, done, nil
}
