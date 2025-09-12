##  Request 

- `Request struct` 
	- `RequestLine RequestLine` // GET /index.html HTTP/1.1
  - `Headers     *headers.Headers` 
  - `State       ParserState string`
    ```go
    - `StateInit    "initialized"` // 1st State
    - `StateDone    "done"`        // Last State
    - `StateHeaders "headers"`     // 2st State, Parsing Headers
    - `StateBody    "body"`        // 3st State, Parsing Body
    ```
  - `Body        string`

- `RequestLine struct` 
  - `Method        string` // GET, POST
  - `RequestTarget string` // /path
  - `HttpVersion   string` // HTTP/1.1, HTTP/2.0

#### Methods / Functions

- `newRequest`- returns Object of `Request`  - Create a new request object
- `validHTTP` - returns `true/false` - Just, If HTTP version and format is right or not
- `parseRequestLine` - returns a `RequestLine`, `length of RequestLine` 
- `RequestFromReader` - returns a whole `Request` - This is too important, Because This reads a data from `Reader`
- `parse` - returns an `int` - This returns a data to `RequestFromReader`, This is a state machine, where every parsing methods are integrated

- For `Request`, Workflow -
  ```markdown
  RequestFromReader -> parse -> parseRequestLine -> parseHeaders
  ```

## Headers

- `Headers struct` 
  - `headers map[string]string` - It's a map. That stores key-value pair 


#### Format 
  `
content-length: 151
connection: close
content-type: text/html
  ` 

#### Method / Functions

- `parseHeaders` - returns a key, value pair of headers - Split a string by `:`, and Checks some cases 
