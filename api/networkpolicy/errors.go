package networkpolicy

import "fmt"

// RawHTTPStatusError represents any response with a 4xx or 5xx status code.
type RawHTTPStatusError struct {
	StatusCode  int
	RawResponse []byte
}

func (r RawHTTPStatusError) Error() string {
	return fmt.Sprintf("Error Code: %d\nRaw Response: %s", r.StatusCode, r.RawResponse)
}

// UnverifiedServerError replaces x509.UnknownAuthorityError when the server
// has SSL but the client is unable to verify it's certificate
type UnverifiedServerError struct {
	URL string
}

func (e UnverifiedServerError) Error() string {
	return "x509: certificate signed by unknown authority"
}

// RequestError represents a generic error encountered while performing the
// HTTP request. This generic error occurs before a HTTP response is obtained.
type RequestError struct {
	Err error
}

func (e RequestError) Error() string {
	return e.Err.Error()
}

// InvalidAuthTokenError is returned when the client has an invalid
// authorization header.
type InvalidAuthTokenError struct {
	Message string
}

func (e InvalidAuthTokenError) Error() string {
	return e.Message
}
