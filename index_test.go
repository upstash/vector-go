package vector

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"io/fs"
	"net/http"
	"os"
	"testing"
)

var (
	namespaces = [...]string{defaultNamespace, "ns"}
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

	for _, ns := range namespaces {
		err := client.Namespace(ns).Reset()
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func newTestClientWithNamespace(ns string) (*Namespace, error) {
	client := NewIndex(
		os.Getenv(UrlEnvProperty),
		os.Getenv(TokenEnvProperty),
	)

	for _, ns := range namespaces {
		err := client.Namespace(ns).Reset()
		if err != nil {
			return nil, err
		}
	}

	return client.Namespace(ns), nil
}

func newEmbeddingTestClient() (*Index, error) {
	client := NewIndex(
		os.Getenv("EMBEDDING_"+UrlEnvProperty),
		os.Getenv("EMBEDDING_"+TokenEnvProperty),
	)

	for _, ns := range namespaces {
		err := client.Namespace(ns).Reset()
		if err != nil {
			return nil, err
		}
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

	for _, ns := range namespaces {
		err := c.Namespace(ns).Reset()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func TestNewClient(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient()
			require.NoError(t, err)

			_, err = client.Info()
			require.NoError(t, err)
		})
	}
}

func TestNewClientWith(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClientWith(&http.Client{})
			require.NoError(t, err)

			_, err = client.Info()
			require.NoError(t, err)
		})
	}
}
