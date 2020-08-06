package models

import "net/http"

// ServerDownResponse is default parameter for response when our server crash.
var ServerDownResponse = Response{
	StatusCode: http.StatusServiceUnavailable,
	Body:       []byte(""),
}

// Response keeps status code and byte body
type Response struct {
	StatusCode int
	Body       []byte
}

// Target used for forcing error on special part of user input.
type Target struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ErrorResponse keeps our error response type.
// target omits when it's empty, because all the responses doesn't have req.body
type ErrorResponse struct {
	Code      string   `json:"code"`
	Message   string   `json:"message"`
	MessageFa string   `json:"messageFa"`
	Target    []Target `json:"target,omitempty"`
}
