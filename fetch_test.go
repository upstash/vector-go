package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
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
					require.Equal(t, v0, vectors[0].Vector)
					require.Equal(t, sv0, vectors[0].SparseVector)
					require.Equal(t, id1, vectors[1].Id)
					require.Nil(t, vectors[1].Metadata) // was upserted with nil metadata
					require.Equal(t, v1, vectors[1].Vector)
					require.Equal(t, sv1, vectors[1].SparseVector)
				})
			})
		}
	}
}
