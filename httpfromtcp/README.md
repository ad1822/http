# HTTP Request Message Format

An HTTP request is a **structured text message** sent by a client to a server. It follows a specific format defined in the HTTP/1.1 specification (RFC 7230–7235). At a high level, it consists of:

1. **Start-line**
2. **Header fields** (also called field-lines)
3. **Blank line** (CRLF)
4. **Optional message body**

---

## 1. Start-Line

The first line of the request. It defines **what** the client wants.

**Format**:

```
<method> <request-target> <HTTP-version>\r\n
```

### Components

* **Method**
  Indicates the desired action (case-sensitive, uppercase).
  Examples:

  * `GET` → retrieve data
  * `POST` → submit data
  * `PUT` → replace a resource
  * `DELETE` → remove a resource
  * `HEAD`, `OPTIONS`, `PATCH`, etc.

* **Request-target**
  The path or resource being requested.

  * Usually: an absolute path (`/index.html`, `/api/v1/users`)
  * May include query (`/search?q=golang`)
  * Special form for proxies: `http://example.com/`

* **HTTP-version**
  Protocol version, e.g. `HTTP/1.0`, `HTTP/1.1`, `HTTP/2.0` (textual in HTTP/1.x only).

**Example**:

```
GET /index.html HTTP/1.1\r\n
```

---

## 2. Header Fields

One or more lines, each representing a **key-value pair** that carries metadata about the request.

**Format**:

```
<Field-Name>: <Field-Value>\r\n
```

### Rules

* Field names are case-insensitive (`Host`, `host`, `HOST` → same meaning).
* Each header ends with `\r\n`.
* No blank lines allowed here (except to mark the end of headers).
* Multiple headers with the same name can be combined or repeated, depending on the field.

### Common Request Headers

* `Host`: required in HTTP/1.1. Example:

  ```
  Host: www.example.com
  ```
* `User-Agent`: identifies the client software.
* `Accept`: media types the client can handle.
* `Content-Length`: size of body in bytes.
* `Content-Type`: type of body (e.g. `application/json`).
* `Authorization`: credentials for authentication.
* `Connection`: options like `keep-alive` or `close`.

---

## 3. Blank Line

After all header fields, a mandatory empty line:

```
\r\n
```

This signals the end of the headers.

---

## 4. Optional Message Body

Not all requests have a body. When present, the body follows the blank line.

* **GET, HEAD, DELETE**: usually no body.
* **POST, PUT, PATCH**: commonly include a body (form data, JSON, files).

The server knows how to interpret the body by looking at headers like:

* `Content-Length` → exact size in bytes.
* `Transfer-Encoding: chunked` → body sent in chunks.
* `Content-Type` → media type (`application/json`, `text/plain`, etc.).

---

## Complete Example

### GET request (no body)

```
GET /search?q=golang HTTP/1.1\r\n
Host: www.example.com\r\n
User-Agent: curl/8.5.0\r\n
Accept: */*\r\n
\r\n
```

### POST request (with body)

```
POST /api/v1/users HTTP/1.1\r\n
Host: api.example.com\r\n
Content-Type: application/json\r\n
Content-Length: 27\r\n
\r\n
{"name":"Ayush","age":22}
```

---

## Summary Checklist

* **Start-line** → method, target, version.
* **Header fields** → one per line, key-value format.
* **Blank line** → separates headers from body.
* **Body** → optional, controlled by headers.

---
