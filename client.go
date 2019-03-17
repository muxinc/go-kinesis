package kinesis

import (
	"net/http"
)

type SignFunc func(Auth, *http.Request) error

// Client is like http.Client, but signs all requests using Auth.
type Client struct {
	// Auth holds the credentials for this client instance
	auth Auth
	// The http client to make requests with. If nil, http.DefaultClient is used.
	client *http.Client

	// The signing function to use for signing the request
	signFunc SignFunc
}

// NewClient creates a new Client that uses the credentials in the specified
// Auth object.
//
// This function assumes the Auth object has been sanely initialized. If you
// wish to infer auth credentials from the environment, refer to NewAuth
func NewClient(auth Auth) *Client {
	return &Client{auth: auth, client: http.DefaultClient, signFunc: Sign}
}

// NewClientWithHTTPClient creates a client with a non-default http client
// ie. a timeout could be set on the HTTP client to timeout if Kinesis doesn't
// response in a timely manner like after the 5 minute mark where the current
// shard iterator expires
func NewClientWithHTTPClient(auth Auth, httpClient *http.Client) *Client {
	return &Client{auth: auth, client: httpClient, signFunc: Sign}
}

// NewRegionClientWithHTTPClient creates a client with a non-default http client with a signing function specified for a region
func NewRegionClientWithHTTPClient(auth Auth, region string, httpClient *http.Client) *Client {
	return &Client{auth: auth, client: httpClient, signFunc: (&Service{Name: "kinesis", Region: region}).Sign}
}

// Do some request, but sign it before sending
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := c.signFunc(c.auth, req)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
