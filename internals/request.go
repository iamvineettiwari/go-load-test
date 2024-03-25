package internals

import "errors"

// Supported method for HTTP(s) request
var (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

// Supported authentication types
var (
	BASIC string = "BASIC"
	TOKEN string = "TOKEN"
)

type RequestAuth struct {
	AuthType string
	Username string
	Password string
	Token    string
}

type Request struct {
	Url     string
	Method  string
	Headers map[string]string
	Body    []byte
	Auth    RequestAuth
}

func NewRequest(url, method string, headers map[string]string, body []byte, auth RequestAuth) *Request {
	return &Request{
		Url:     url,
		Method:  method,
		Headers: headers,
		Body:    body,
		Auth:    auth,
	}
}

func NewRequestAuth(basicAuth, tokenAuth bool, username, password, token string) (RequestAuth, error) {
	reqAuth := RequestAuth{}

	if basicAuth && (username == "" || password == "") {
		return reqAuth, errors.New("username and password is required for basic auth")
	}

	if tokenAuth && token == "" {
		return reqAuth, errors.New("token is required for token auth")
	}

	if basicAuth {
		reqAuth.AuthType = BASIC
		reqAuth.Username = username
		reqAuth.Password = password
		return reqAuth, nil
	}

	reqAuth.AuthType = TOKEN
	reqAuth.Token = token

	return reqAuth, nil
}
