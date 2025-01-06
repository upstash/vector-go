package vector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReset(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeDense, ns)
			require.NoError(t, err)

			id := randomString()
			err = client.Upsert(Upsert{
				Id:     id,
				Vector: []float32{0, 1},
			})
			require.NoError(t, err)

			err = client.Reset()
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				info, err := client.Info()
				require.NoError(t, err)
				return info.VectorCount == 0
			}, 10*time.Second, 1*time.Second)
		})
	}
}
