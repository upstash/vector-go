package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
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

	t.Run("single", func(t *testing.T) {
		vectors, err := client.Fetch(Fetch{
			Ids: []string{id0},
		})
		require.NoError(t, err)

		require.Equal(t, 1, len(vectors))
		require.Equal(t, id0, vectors[0].Id)
	})

	t.Run("many", func(t *testing.T) {
		vectors, err := client.Fetch(Fetch{
			Ids: []string{id0, id1},
		})
		require.NoError(t, err)

		require.Equal(t, 2, len(vectors))
		require.Equal(t, id0, vectors[0].Id)
		require.Equal(t, id1, vectors[1].Id)
	})

	t.Run("non existing id", func(t *testing.T) {
		vectors, err := client.Fetch(Fetch{
			Ids: []string{randomString()},
		})
		require.NoError(t, err)

		require.Equal(t, 1, len(vectors))
		require.Equal(t, "", vectors[0].Id)
	})

	t.Run("with metadata and vectors", func(t *testing.T) {
		vectors, err := client.Fetch(Fetch{
			Ids:             []string{id0, id1},
			IncludeMetadata: true,
			IncludeVectors:  true,
		})
		require.NoError(t, err)

		require.Equal(t, 2, len(vectors))
		require.Equal(t, id0, vectors[0].Id)
		require.Equal(t, map[string]any{"foo": "bar"}, vectors[0].Metadata)
		require.Equal(t, []float32{0, 1}, vectors[0].Vector)
		require.Equal(t, id1, vectors[1].Id)
		require.Nil(t, vectors[1].Metadata) // was upserted with nil metadata
		require.Equal(t, []float32{5, 10}, vectors[1].Vector)
	})
}
