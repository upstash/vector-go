package vector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	client, err := newTestClient()
	require.NoError(t, err)

	id0 := randomString()
	id1 := randomString()
	err = client.UpsertMany([]Upsert{
		{
			Id:       id0,
			Vector:   []float32{0, 1},
			Metadata: map[string]any{"foo": "bar"},
		},
		{
			Id:     id1,
			Vector: []float32{5, 10},
		},
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		info, err := client.Info()
		require.NoError(t, err)
		return info.PendingVectorCount == 0
	}, 10*time.Second, 1*time.Second)

	t.Run("score", func(t *testing.T) {
		scores, err := client.Query(Query{
			Vector: []float32{0, 1},
			TopK:   1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(scores))
		require.Equal(t, id0, scores[0].Id)
		require.Equal(t, float32(1.0), scores[0].Score)
	})

	t.Run("with metadata and vectors", func(t *testing.T) {
		scores, err := client.Query(Query{
			Vector:          []float32{0, 1},
			TopK:            1,
			IncludeMetadata: true,
			IncludeVectors:  true,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(scores))
		require.Equal(t, id0, scores[0].Id)
		require.Equal(t, float32(1.0), scores[0].Score)
		require.Equal(t, map[string]any{"foo": "bar"}, scores[0].Metadata)
		require.Equal(t, []float32{0, 1}, scores[0].Vector)
	})
}
