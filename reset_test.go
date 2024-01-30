package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReset(t *testing.T) {
	client, err := newTestClient()
	require.NoError(t, err)

	id := randomString()
	err = client.Upsert(Upsert{
		Id:     id,
		Vector: []float32{0, 1},
	})
	require.NoError(t, err)

	err = client.Reset()
	require.NoError(t, err)

	info, err := client.Info()
	require.NoError(t, err)
	require.Equal(t, 0, info.VectorCount)
}
