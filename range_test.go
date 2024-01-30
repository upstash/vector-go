package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {
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

	t.Run("range all", func(t *testing.T) {
		vectors, err := client.Range(Range{
			Cursor: "",
			Limit:  2,
		})
		require.NoError(t, err)

		require.Equal(t, 2, len(vectors.Vectors))
		require.Equal(t, id0, vectors.Vectors[0].Id)
		require.Equal(t, id1, vectors.Vectors[1].Id)
		require.Equal(t, "", vectors.NextCursor)
	})

	t.Run("range part by part", func(t *testing.T) {
		vectors, err := client.Range(Range{
			Cursor: "",
			Limit:  1,
		})
		require.NoError(t, err)

		require.Equal(t, 1, len(vectors.Vectors))
		require.Equal(t, id0, vectors.Vectors[0].Id)
		require.Equal(t, "1", vectors.NextCursor)

		vectors, err = client.Range(Range{
			Cursor: "1",
		})
		require.NoError(t, err)

		require.Equal(t, 1, len(vectors.Vectors))
		require.Equal(t, id1, vectors.Vectors[0].Id)
		require.Equal(t, "", vectors.NextCursor)
	})

	t.Run("with metadata and vectors", func(t *testing.T) {
		vectors, err := client.Range(Range{
			Limit:           2,
			IncludeMetadata: true,
			IncludeVectors:  true,
		})
		require.NoError(t, err)

		require.Equal(t, 2, len(vectors.Vectors))
		require.Equal(t, id0, vectors.Vectors[0].Id)
		require.Equal(t, map[string]any{"foo": "bar"}, vectors.Vectors[0].Metadata)
		require.Equal(t, []float32{0, 1}, vectors.Vectors[0].Vector)
		require.Equal(t, id1, vectors.Vectors[1].Id)
		require.Nil(t, vectors.Vectors[1].Metadata) // was upserted with nil metadata
		require.Equal(t, []float32{5, 10}, vectors.Vectors[1].Vector)
		require.Equal(t, "", vectors.NextCursor)
	})
}
