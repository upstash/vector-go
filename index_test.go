package vector

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(err)
	}
}

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
