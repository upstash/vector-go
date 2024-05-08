package vector

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNamespace(t *testing.T) {

	t.Run("list namespaces", func(t *testing.T) {
		client, err := newTestClient()
		require.NoError(t, err)

		for _, ns := range namespaces {
			createNamespace(t, client, ns)
		}

		ns, err := client.ListNamespaces()
		require.NoError(t, err)
		require.Exactly(t, []string{"", "ns"}, ns)
	})

	t.Run("delete namespaces", func(t *testing.T) {
		client, err := newTestClient()
		require.NoError(t, err)

		for _, ns := range namespaces {
			createNamespace(t, client, ns)
		}

		for _, ns := range namespaces {
			if ns == "" {
				continue
			}

			err := client.DeleteNamespace(ns)
			require.NoError(t, err)
		}

		info, err := client.Info()
		require.NoError(t, err)
		require.Exactly(t, 1, len(info.NamespaceInfo))
	})
}

func createNamespace(t *testing.T, client *Index, ns string) {
	if ns == "" {
		return
	}

	client = client.Namespace(ns)

	err := client.Upsert(Upsert{
		Vector: []float32{0.1, 0.1},
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		info, err := client.Info()
		require.NoError(t, err)
		return info.PendingVectorCount == 0
	}, 10*time.Second, 1*time.Second)

	err = client.Reset()
	require.NoError(t, err)
}
