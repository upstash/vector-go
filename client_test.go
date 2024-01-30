package vector

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestClient() (*Client, error) {
	client := NewClient(
		os.Getenv("UPSTASH_VECTOR_REST_URL"),
		os.Getenv("UPSTASH_VECTOR_REST_TOKEN"),
	)

	err := client.Reset()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newTestClientWith(client *http.Client) (*Client, error) {
	opts := Options{
		Url:    os.Getenv("UPSTASH_VECTOR_REST_URL"),
		Token:  os.Getenv("UPSTASH_VECTOR_REST_TOKEN"),
		Client: client,
	}

	c := NewClientWith(opts)

	err := c.Reset()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func TestNewClient(t *testing.T) {
	client, err := newTestClient()
	require.NoError(t, err)

	_, err = client.Info()
	require.NoError(t, err)
}

func TestNewClientWith(t *testing.T) {
	client, err := newTestClientWith(&http.Client{})
	require.NoError(t, err)

	_, err = client.Info()
	require.NoError(t, err)
}
