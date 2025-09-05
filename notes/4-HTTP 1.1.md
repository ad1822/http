# HTTP/1.1


- Stateless application level req/res protocol that uses extensible semantics and self-descriptive messages for flexible interaction with network-based hyprtext information systems
- Defined By :
  - HTTP semantics
  - HTTP caching

- HTTP/1.1 clients and servers communicate by sending messages. See Section 3 of [HTTP] for the general terminology and core concepts of HTTP.

## Message Format

- An HTTP/1.1 message consists of a start-line followed by a CRLF and a sequence of octets in a format similar to the Internet Message Format [RFC5322]: zero or more header field lines (collectively referred to as the "headers" or the "header section"), an empty line indicating the end of the header section, and an optional message body.

```
  HTTP-message   = start-line CRLF
                   *( field-line CRLF )
                   CRLF
                   [ message-body ]
```
- A message can be either a request from client or a response from server.
- The two types of messages differ only in the start-line, which is either a request-line (from requests) or a status-line (from responses)


```
start-line = request-line / status-line
```
- In theory, a client could receive requests and a server could receive responses, distinguishing them by their different start-line formats. In practice, servers are implemented to only expect a request (a response is interpreted as an unknown or invalid request method), and clients are implemented to only expect a
- HTTP makes use of some protocol elements similar to the Multipurpose Internet Mail Extensions (MIME)

# Message Parsing

- Read the start-line into a structure representation
- Accumulate header field lines into a loopup keyed by field name, until the blank line signaling end-of-headers
- Determine if a message body exists using header information
- If a body is expected, read exactly the declared number of bytes (via `Content-Length`, `Trasfer-Encoding`  ), or until connection closure
- **String-based Parsers** : Only safe for processing content after the message has been structurally parsed—e.g., inside header values after field boundaries are isolated.

