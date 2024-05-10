package vector

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	client, err := newTestClient()
	require.NoError(t, err)

	for _, ns := range namespaces {
		createNamespace(t, client, ns)
	}

	info, err := client.Info()
	require.NoError(t, err)
	require.Equal(t, info.VectorCount, 0)
	require.Equal(t, 2, info.Dimension)
	require.Equal(t, "COSINE", info.SimilarityFunction)
	require.Equal(t, len(namespaces), len(info.Namespaces))
	for _, ns := range namespaces {
		require.Contains(t, info.Namespaces, ns)
		require.Equal(t, 0, info.Namespaces[ns].VectorCount)
		require.Equal(t, 0, info.Namespaces[ns].PendingVectorCount)
	}

	for _, ns := range namespaces {
		err = client.Namespace(ns).Upsert(Upsert{
			Id:     randomString(),
			Vector: []float32{0, 1},
		})
		require.NoError(t, err)
	}

	require.Eventually(t, func() bool {
		info, err := client.Info()
		require.NoError(t, err)
		return info.VectorCount == len(namespaces)
	}, 10*time.Second, 1*time.Second)
}
