package vector

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestClient() (*Index, error) {
	client := NewIndex(
		os.Getenv(UrlEnvProperty),
		os.Getenv(TokenEnvProperty),
	)

	err := client.Reset()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newTestClientWith(client *http.Client) (*Index, error) {
	opts := Options{
		Url:    os.Getenv(UrlEnvProperty),
		Token:  os.Getenv(TokenEnvProperty),
		Client: client,
	}

	c := NewIndexWith(opts)

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
