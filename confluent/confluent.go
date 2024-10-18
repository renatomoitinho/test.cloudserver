package confluent

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ClientRequestError struct {
	Err     error
	Status  int
	Message string
}

func (c *ClientRequestError) Error() string {
	return c.Err.Error()
}

func (c *ClientRequestError) Unmarshal(ref any) error {
	return json.Unmarshal([]byte(c.Message), ref)
}

type ClientRequest struct {
	client  *http.Client
	request *http.Request
}

func (c *ClientRequest) Get(ops ...func(*ClientRequest)) ([]byte, error) {
	var (
		err  error
		resp *http.Response
	)

	for _, op := range ops {
		op(c)
	}

	resp, err = c.client.Do(c.request)
	if err != nil {
		return nil, &ClientRequestError{
			Err: err,
		}
	}
	defer resp.Body.Close()

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ClientRequestError{
			Err:     err,
			Status:  resp.StatusCode,
			Message: string(body),
		}
	}

	if resp.StatusCode/100 != 2 {
		return nil, &ClientRequestError{
			Err:     errors.New("request status code is not 2xx"),
			Status:  resp.StatusCode,
			Message: string(body),
		}
	}

	return body, nil
}

func WithClientInsecure() func(*ClientRequest) {
	return func(c *ClientRequest) {
		c.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

func WithBasicAuthorization(username, password string) func(*ClientRequest) {
	return func(r *ClientRequest) {
		auth := username + ":" + password
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		r.request.Header.Add("Authorization", authHeader)
	}
}

func WithUrl(urlstr string) func(*ClientRequest) {
	return func(r *ClientRequest) {
		u, _ := url.Parse(urlstr)
		r.request.URL = u
	}
}

func WithHeader(key, value string) func(*ClientRequest) {
	return func(r *ClientRequest) {
		r.request.Header.Add(key, value)
	}
}

func WithParam(key, value string) func(*ClientRequest) {
	return func(r *ClientRequest) {
		q := r.request.URL.Query()
		q.Add(key, value)
		r.request.URL.RawQuery = q.Encode()
	}
}

func Get[T any](url string, ops ...func(*ClientRequest)) (*T, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c := &ClientRequest{
		client:  http.DefaultClient,
		request: req,
	}

	body, err := c.Get(ops...)
	if err != nil {
		return nil, err
	}
	var res T
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", string(body))
	}
	return &res, nil
}

func HasError(err error) *ClientRequestError {
	if err == nil {
		return nil
	}
	var e *ClientRequestError
	if errors.As(err, &e) {
		return e
	}
	return &ClientRequestError{
		Err: err,
	}
}
