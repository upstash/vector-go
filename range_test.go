package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {
	for _, ns := range namespaces {
		for _, tcType := range testClientTypes {
			t.Run("namespace_"+ns+"_index_type_"+string(tcType), func(t *testing.T) {
				client, err := newTestClient(tcType, ns)
				require.NoError(t, err)

				id0 := randomString()
				v0, sv0 := randomVectors(tcType)

				id1 := randomString()
				v1, sv1 := randomVectors(tcType)

				err = client.UpsertMany([]Upsert{
					{
						Id:           id0,
						Vector:       v0,
						SparseVector: sv0,
						Metadata:     map[string]any{"foo": "bar"},
					},
					{
						Id:           id1,
						Vector:       v1,
						SparseVector: sv1,
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
					require.Equal(t, v0, vectors.Vectors[0].Vector)
					require.Equal(t, sv0, vectors.Vectors[0].SparseVector)
					require.Equal(t, id1, vectors.Vectors[1].Id)
					require.Nil(t, vectors.Vectors[1].Metadata) // was upserted with nil metadata
					require.Equal(t, v1, vectors.Vectors[1].Vector)
					require.Equal(t, sv1, vectors.Vectors[1].SparseVector)
					require.Equal(t, "", vectors.NextCursor)
				})
			})
		}
	}
}
