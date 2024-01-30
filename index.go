package vector

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const (
	UrlEnvProperty   = "UPSTASH_VECTOR_REST_URL"
	TokenEnvProperty = "UPSTASH_VECTOR_REST_TOKEN"
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
	if o.Url == "" {
		panic("Missing Upstash Vector URL")
	}
	if o.Token == "" {
		panic("Missing Upstash Vector Token")
	}
}

// NewIndex returns an index client to be used with Upstash Vector
// with the given url and token.
func NewIndex(url string, token string) *Index {
	return NewIndexWith(Options{
		Url:   url,
		Token: token,
	})
}

// NewIndexFromEnv returns an index client to be used with Upstash Vector
// by reading URL and token from the environment variables.
func NewIndexFromEnv() *Index {
	return NewIndexWith(Options{
		Url:   os.Getenv(UrlEnvProperty),
		Token: os.Getenv(TokenEnvProperty),
	})
}

// NewIndexWith returns an index client to be used with Upstash Vector
// with the given options.
func NewIndexWith(options Options) *Index {
	options.init()
	return &Index{
		url:    options.Url,
		token:  options.Token,
		client: options.Client,
	}
}

// Index is a client for Upstash Vector index.
type Index struct {
	url    string
	token  string
	client *http.Client
}

func (ix *Index) sendJson(path string, obj any) (data []byte, err error) {
	if data, err = json.Marshal(obj); err != nil {
		return
	}
	return ix.sendBytes(path, data)
}

func (ix *Index) sendBytes(path string, obj []byte) (data []byte, err error) {
	return ix.send(path, bytes.NewReader(obj))
}

func (ix *Index) send(path string, r io.Reader) (data []byte, err error) {
	request, err := http.NewRequest("POST", ix.url+path, r)
	if err != nil {
		return
	}
	request.Header.Add("Authorization", "Bearer "+ix.token)
	response, err := ix.client.Do(request)
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
