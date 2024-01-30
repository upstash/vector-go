package vector

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Options struct {
	// URL of the Upstash Vector index.
	Url string

	// Token of the Upstash Vector index.
	Token string

	// The HTTP client to use for requests.
	Client *http.Client
}

func (o *Options) init() {
	if o.Client == nil {
		o.Client = http.DefaultClient
	}
}

// NewClient returns a client to be used with Upstash Vector
// with the given url and token.
func NewClient(url string, token string) *Client {
	return NewClientWith(Options{
		Url:   url,
		Token: token,
	})
}

// NewClient returns a client to be used with Upstash Vector
// with the given options.
func NewClientWith(options Options) *Client {
	options.init()
	return &Client{
		url:    options.Url,
		token:  options.Token,
		client: options.Client,
	}
}

type Client struct {
	url    string
	token  string
	client *http.Client
}

func (c *Client) sendJson(path string, obj any) (data []byte, err error) {
	if data, err = json.Marshal(obj); err != nil {
		return
	}
	return c.sendBytes(path, data)
}

func (c *Client) sendString(path string, obj string) (data []byte, err error) {
	return c.send(path, strings.NewReader(obj))
}

func (c *Client) sendBytes(path string, obj []byte) (data []byte, err error) {
	return c.send(path, bytes.NewReader(obj))
}

func (c *Client) send(path string, r io.Reader) (data []byte, err error) {
	request, err := http.NewRequest("POST", c.url+path, r)
	if err != nil {
		return
	}
	request.Header.Add("Authorization", "Bearer "+c.token)
	response, err := c.client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	data, err = io.ReadAll(response.Body)
	return
}

func parseResponse[T any](data []byte) (t T, err error) {
	var result response[T]
	if err = json.Unmarshal(data, &result); err != nil {
		return
	}
	t = result.Result
	if result.Error != "" {
		err = errors.New(result.Error)
	}
	return
}
