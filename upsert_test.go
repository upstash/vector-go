package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpsert(t *testing.T) {
	client, err := newTestClient()
	require.NoError(t, err)

	t.Run("single", func(t *testing.T) {
		id := randomString()
		err := client.Upsert(Upsert{
			Id:     id,
			Vector: []float32{0, 1},
		})
		require.NoError(t, err)

		vectors, err := client.Fetch(Fetch{
			Ids: []string{id},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(vectors))
		require.Equal(t, id, vectors[0].Id)
	})

	t.Run("many", func(t *testing.T) {
		id0 := randomString()
		id1 := randomString()
		err = client.UpsertMany([]Upsert{
			{
				Id:     id0,
				Vector: []float32{0, 1},
			},
			{
				Id:     id1,
				Vector: []float32{5, 10},
			},
		})
		require.NoError(t, err)

		vectors, err := client.Fetch(Fetch{
			Ids: []string{id0, id1},
		})
		require.NoError(t, err)
		require.Equal(t, 2, len(vectors))
		require.Equal(t, id0, vectors[0].Id)
		require.Equal(t, id1, vectors[1].Id)
	})

	t.Run("with metadata", func(t *testing.T) {
		id := randomString()
		err := client.Upsert(Upsert{
			Id:       id,
			Vector:   []float32{0, 1},
			Metadata: map[string]any{"foo": "bar"},
		})
		require.NoError(t, err)

		vectors, err := client.Fetch(Fetch{
			Ids:             []string{id},
			IncludeMetadata: true,
			IncludeVectors:  true,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(vectors))
		require.Equal(t, id, vectors[0].Id)
		require.Equal(t, map[string]any{"foo": "bar"}, vectors[0].Metadata)
		require.Equal(t, []float32{0, 1}, vectors[0].Vector)
	})
}
