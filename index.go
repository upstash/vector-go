package vector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

const (
	UrlEnvProperty   = "UPSTASH_VECTOR_REST_URL"
	TokenEnvProperty = "UPSTASH_VECTOR_REST_TOKEN"
	defaultNamespace = ""
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
	index := &Index{
		url:    options.Url,
		token:  options.Token,
		client: options.Client,
	}
	index.generateHeaders()
	return index
}

// Index is a client for Upstash Vector index.
type Index struct {
	url     string
	token   string
	client  *http.Client
	headers http.Header
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
	request, err := http.NewRequest(http.MethodPost, ix.url+path, r)
	if err != nil {
		return
	}
	request.Header = ix.headers
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

func (ix *Index) generateHeaders() {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+ix.token)
	headers.Add("Upstash-Telemetry-Runtime", fmt.Sprintf("vector-go@%s", runtime.Version()))
	var platform string
	if os.Getenv("VERCEL") != "" {
		platform = "vercel"
	} else if os.Getenv("AWS_REGION") != "" {
		platform = "aws"
	} else {
		platform = "unknown"
	}
	headers.Add("Upstash-Telemetry-Platform", platform)
	ix.headers = headers
}

func buildPath(path string, ns string) string {
	if ns == "" {
		return path
	}
	return path + "/" + ns
}
