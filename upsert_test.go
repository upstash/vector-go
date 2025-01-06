package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpsert(t *testing.T) {
	for _, ns := range namespaces {
		for _, tcType := range testClientTypes {
			t.Run("namespace_"+ns+"_index_type_"+string(tcType), func(t *testing.T) {
				client, err := newTestClient(tcType, ns)
				require.NoError(t, err)

				t.Run("single", func(t *testing.T) {
					id := randomString()
					v0, sv0 := randomVectors(tcType)

					err := client.Upsert(Upsert{
						Id:           id,
						Vector:       v0,
						SparseVector: sv0,
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
					v0, sv0 := randomVectors(tcType)

					id1 := randomString()
					v1, sv1 := randomVectors(tcType)

					err = client.UpsertMany([]Upsert{
						{
							Id:           id0,
							Vector:       v0,
							SparseVector: sv0,
						},
						{
							Id:           id1,
							Vector:       v1,
							SparseVector: sv1,
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
					v0, sv0 := randomVectors(tcType)

					err := client.Upsert(Upsert{
						Id:           id,
						Vector:       v0,
						SparseVector: sv0,
						Metadata:     map[string]any{"foo": "bar"},
						Data:         "some data",
					})
					require.NoError(t, err)

					vectors, err := client.Fetch(Fetch{
						Ids:             []string{id},
						IncludeMetadata: true,
						IncludeVectors:  true,
						IncludeData:     true,
					})
					require.NoError(t, err)
					require.Equal(t, 1, len(vectors))
					require.Equal(t, id, vectors[0].Id)
					require.Equal(t, map[string]any{"foo": "bar"}, vectors[0].Metadata)
					require.Equal(t, v0, vectors[0].Vector)
					require.Equal(t, sv0, vectors[0].SparseVector)
					require.Equal(t, "some data", vectors[0].Data)
				})
			})
		}
	}
}
