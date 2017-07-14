package networkpolicy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// NetworkPolicyConnection represents the connection to NetworkPolicy
type NetworkPolicyConnection struct {
	HTTPClient *http.Client
}

// NewConnection returns a pointer to a new NetworkPolicy Connection
func NewConnection(skipSSLValidation bool, dialTimeout time.Duration) *NetworkPolicyConnection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLValidation,
		},
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			Timeout:   dialTimeout,
		}).DialContext,
	}

	return &NetworkPolicyConnection{
		HTTPClient: &http.Client{Transport: tr},
	}
}

func (connection *NetworkPolicyConnection) Make(request *http.Request, passedResponse *Response) error {
	// In case this function is called from a retry, passedResponse may already
	// be populated with a previous response. We reset in case there's an HTTP
	// error and we don't repopulate it in populateResponse.
	passedResponse.reset()
	fmt.Println(request.URL)
	fmt.Printf("%s\n", request.Body)

	response, err := connection.HTTPClient.Do(request)
	if err != nil {
		return connection.processRequestErrors(request, err)
	}

	return connection.populateResponse(response, passedResponse)
}

func (connection *NetworkPolicyConnection) processRequestErrors(request *http.Request, err error) error {
	switch e := err.(type) {
	case *url.Error:
		if _, ok := e.Err.(x509.UnknownAuthorityError); ok {
			return UnverifiedServerError{
				URL: request.URL.String(),
			}
		}
		return RequestError{Err: e}
	default:
		return err
	}
}

func (connection *NetworkPolicyConnection) populateResponse(response *http.Response, passedResponse *Response) error {
	passedResponse.HTTPResponse = response

	rawBytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return err
	}
	passedResponse.RawResponse = rawBytes

	err = connection.handleStatusCodes(response, passedResponse)
	if err != nil {
		return err
	}

	if passedResponse.Result != nil {
		decoder := json.NewDecoder(bytes.NewBuffer(passedResponse.RawResponse))
		decoder.UseNumber()
		err = decoder.Decode(passedResponse.Result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*NetworkPolicyConnection) handleStatusCodes(response *http.Response, passedResponse *Response) error {
	if response.StatusCode >= 400 {
		return RawHTTPStatusError{
			StatusCode:  response.StatusCode,
			RawResponse: passedResponse.RawResponse,
		}
	}

	return nil
}
