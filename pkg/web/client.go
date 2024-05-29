package web

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

var (
	maxBufferSize = 1000 * 512 // 512kb
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
	respDump := []byte{}
	if err == nil {
		respDump, _ = httputil.DumpResponse(resp, c.dumpBody(resp.Header.Get("Content-Type")))
	} else {
		respDump = []byte("request failed with error: " + err.Error())
	}

	c.log.Logf("[DEBUG] http client dump:\n\nrequest: %s\n\nresponse: %s\n\n", string(reqDump), string(respDump))
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
	Stream(part any) (<-chan any, error)
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

// Stream - предназначен для отправки HTTP-запросов и получения потоковых ответов в формате JSON.
// Метод создает новый HTTP-запрос с заданным путем, методами запроса, заголовками и телом,
// отправляет его и обрабатывает потоковый ответ.
//
// #### Входные параметры:
// - `part`: Переменная, в которую будут декодироваться части потокового ответа.
//
// #### Возвращаемые значения:
// - `chan<- any`: Канал, в который будут отправляться декодированные части ответа.
// - `error`: Ошибка, если операция завершилась неуспешно.
//
// #### Пример использования:
// ```go
// client := NewClient(logger, baseURL)
// jsonReq := client.JSON("/path/to/resource").SetMethod(http.MethodGet)
// stream, err := jsonReq.Stream(part)
//
//	if err != nil {
//	    log.Fatalf("Failed to stream: %v", err)
//	}
//
//	for response := range stream {
//	    // Обработка каждого фрагмента ответа...
//	}
//
// ```
func (r *jsonRequest) Stream(part any) (<-chan any, error) {
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
			return nil, err
		}
	}

	req, err := http.NewRequest(r.method, url.String(), buffer)
	if err != nil {
		r.client.log.Logf("[ERROR] failed to create request: %v", err)
		return nil, err
	}

	for key, values := range r.header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		r.client.log.Logf("[ERROR] failed to send request: %v", err)
		return nil, err
	}

	if resp.StatusCode >= 300 {
		r.client.log.Logf("[ERROR] stream request failed with status code: %v", resp.StatusCode)
		return nil, &RequestError{
			StatusCode: resp.StatusCode,
			Body: func(b []byte, err error) []byte {
				if err != nil {
					return []byte(err.Error())
				}

				return b
			}(io.ReadAll(resp.Body)),
		}
	}

	out := make(chan any)

	go func() {
		scanner := bufio.NewScanner(resp.Body)
		scanBuf := make([]byte, 0, maxBufferSize)
		scanner.Buffer(scanBuf, maxBufferSize)

		for scanner.Scan() {
			bts := scanner.Bytes()

			err = json.NewDecoder(bytes.NewReader(bts)).Decode(&part)
			if err != nil {
				r.client.log.Logf("[ERROR] failed to decode response body: %v", err)
				continue
			}

			out <- part
		}

		close(out)
	}()

	return out, nil
}
