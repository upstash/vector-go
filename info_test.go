package vector

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInfo(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient()
			require.NoError(t, err)

			err = client.Upsert(Upsert{
				Id:     randomString(),
				Vector: []float32{0, 1},
			})
			require.NoError(t, err)

			info, err := client.Info()
			require.NoError(t, err)
			require.Greater(t, info.VectorCount, 0)
			require.Equal(t, 2, info.Dimension)
			require.Equal(t, "COSINE", info.SimilarityFunction)
		})
	}
}
