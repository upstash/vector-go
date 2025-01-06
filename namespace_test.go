package vector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNamespace(t *testing.T) {
	t.Run("list namespaces", func(t *testing.T) {
		index := NewIndexFromEnv()

		for _, ns := range namespaces {
			createNamespace(t, index, ns)
		}

		ns, err := index.ListNamespaces()
		require.NoError(t, err)
		require.Exactly(t, []string{"", "ns"}, ns)
	})

	t.Run("delete namespaces", func(t *testing.T) {
		index := NewIndexFromEnv()

		for _, ns := range namespaces {
			createNamespace(t, index, ns)
		}

		for _, ns := range namespaces {
			if ns == "" {
				continue
			}

			err := index.Namespace(ns).DeleteNamespace()
			require.NoError(t, err)
		}

		info, err := index.Info()
		require.NoError(t, err)
		require.Exactly(t, 1, len(info.Namespaces))
	})
}

func createNamespace(t *testing.T, index *Index, ns string) {
	if ns == "" {
		err := index.Reset()
		require.NoError(t, err)
		return
	}

	namespace := index.Namespace(ns)
	err := namespace.Reset()
	require.NoError(t, err)

	err = namespace.Upsert(Upsert{
		Vector: []float32{0.1, 0.1},
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		info, err := index.Info()
		require.NoError(t, err)
		return info.PendingVectorCount == 0
	}, 10*time.Second, 1*time.Second)

	err = namespace.Reset()
	require.NoError(t, err)
}
