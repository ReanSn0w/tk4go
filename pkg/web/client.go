package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// MSRK: - Error
// Конфтрукция для хранения ошибки полученной при ответе от сервера

type RequestError struct {
	StatusCode int
	Body       []byte
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("requst error (code: %v): %s", r.StatusCode, string(r.Body))
}

func (r *RequestError) Decode(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// MARK: - JSON
// Конструкция для отправки запросов через с обработкой JSON

func (c *Client) JSON(path string) JSON {
	return &jsonRequest{
		client: c,
		path:   path,
		header: make(http.Header),
		query:  make(url.Values),
		method: http.MethodGet,
		body:   nil,
	}
}

type JSON interface {
	SetMethod(method string) JSON
	SetHeader(key, value string) JSON
	SetQuery(key, value string) JSON
	SetBody(body interface{}) JSON
	Do(obj any) error
}

type jsonRequest struct {
	client *Client

	method string
	path   string
	header http.Header
	query  url.Values
	body   interface{}
}

func (r *jsonRequest) SetMethod(method string) JSON {
	r.method = method
	return r
}

func (r *jsonRequest) SetHeader(key, value string) JSON {
	r.header.Set(key, value)
	return r
}

func (r *jsonRequest) SetQuery(key, value string) JSON {
	r.header.Set(key, value)
	return r
}

func (r *jsonRequest) SetBody(body interface{}) JSON {
	r.body = body
	return r
}

func (r *jsonRequest) Do(obj any) error {
	url := &url.URL{
		Scheme:   r.client.base.Scheme,
		Host:     r.client.base.Host,
		User:     r.client.base.User,
		Path:     r.path,
		RawQuery: r.query.Encode(),
	}

	buffer := new(bytes.Buffer)
	if r.body != nil {
		err := json.NewEncoder(buffer).Encode(r.body)
		if err != nil {
			r.client.log.Logf("[ERROR] failed to encode request body: %v", err)
			return err
		}
	}

	req, err := http.NewRequest(r.method, url.String(), buffer)
	if err != nil {
		r.client.log.Logf("[ERROR] failed to create request: %v", err)
		return err
	}

	for key, values := range r.header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		r.client.log.Logf("[ERROR] failed to send request: %v", err)
		return err
	}

	if resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			r.client.log.Logf("[ERROR] failed to read response body: %v", err)
			return err
		}

		r.client.log.Logf("[ERROR] request failed with status code: %v", resp.StatusCode)
		return &RequestError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	err = json.NewDecoder(resp.Body).Decode(obj)
	if err != nil {
		r.client.log.Logf("[ERROR] failed to decode response body: %v", err)
		return err
	}

	return nil
}
