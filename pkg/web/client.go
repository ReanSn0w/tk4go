package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

// NewClient creates a new client instance
func NewClient(logger tools.Logger, baseURL string, prepare ...PrepareClientFn) (client *Client, err error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client = &Client{
		base: base,
		log:  logger,
		cl:   prepareClient(prepare).prepare(&http.Client{}),
	}

	return
}

type (
	PrepareClientFn  func(*http.Client) *http.Client
	PrepareRequestFn func(*http.Request) *http.Request
	prepareClient    []PrepareClientFn
	prepareRequest   []PrepareRequestFn

	Client struct {
		base *url.URL
		log  tools.Logger
		cl   *http.Client
	}
)

func (p prepareClient) prepare(client *http.Client) *http.Client {
	for _, fn := range p {
		client = fn(client)
	}

	return client
}

func (p prepareRequest) prepare(req *http.Request) *http.Request {
	for _, fn := range p {
		req = fn(req)
	}

	return req
}

// PrepareClient returns a new client with the given prepare functions
func (c *Client) Json(method string, path string, query url.Values, reqBody, respBody any) error {
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	err := enc.Encode(reqBody)
	if err != nil {
		c.log.Logf("[ERROR] failed to encode request body: %v", err)
		return err
	}

	req, err := http.NewRequest(method, (&url.URL{
		Scheme:   c.base.Scheme,
		Host:     c.base.Host,
		User:     c.base.User,
		Path:     path,
		RawQuery: query.Encode(),
	}).String(), buffer)
	if err != nil {
		c.log.Logf("[ERROR] failed to create request: %v", err)
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		c.log.Logf("[ERROR] failed to send request: %v", err)
		return err
	}

	if resp.StatusCode > 300 {
		c.log.Logf("[ERROR] request failed with status code: %v", resp.StatusCode)

		respErrorData := Response[tools.ErrorsMap]{}
		err := json.NewDecoder(resp.Body).Decode(&respErrorData)
		if err != nil {
			c.log.Logf("[ERROR] failed to decode failure response body: %v", err)
			return err
		}

		return respErrorData.isError()
	}

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		c.log.Logf("[ERROR] failed to decode success response body: %v", err)
		return err
	}

	return nil
}

// Do sends the request and returns the response
func (c *Client) Do(req *http.Request, prepare ...PrepareRequestFn) (*http.Response, error) {
	req = prepareRequest(prepare).prepare(req)
	reqDump, _ := httputil.DumpRequest(req, c.dumpBody(req.Header.Get("Content-Type")))

	resp, err := c.cl.Do(req)
	respDump, _ := httputil.DumpResponse(resp, c.dumpBody(resp.Header.Get("Content-Type")))

	c.log.Logf("[DEBUG] http client dump:\n\nrequest: %s\nresponse: %s\n", string(reqDump), string(respDump))
	return resp, err
}

func (c *Client) dumpBody(contentType string) bool {
	switch contentType {
	case "form/urlencoded":
		return true
	case "application/json":
		return true
	case "application/xml":
		return true
	case "text/xml":
		return true
	case "text/plain":
		return true
	default:
		return false
	}
}
