package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	for _, ns := range namespaces {
		for _, tcType := range testClientTypes {
			t.Run("namespace_"+ns+"_index_type_"+string(tcType), func(t *testing.T) {
				client, err := newTestClient(tcType, ns)
				require.NoError(t, err)

				t.Run("not found", func(t *testing.T) {
					id := randomString()
					ok, err := client.Update(Update{
						Id:   id,
						Data: "new-data",
					})
					require.NoError(t, err)
					require.False(t, ok)
				})

				t.Run("data", func(t *testing.T) {
					id := randomString()
					v, sv := randomVectors(tcType)

					err := client.Upsert(Upsert{
						Id:           id,
						Vector:       v,
						SparseVector: sv,
						Data:         "old-data",
					})
					require.NoError(t, err)

					ok, err := client.Update(Update{
						Id:   id,
						Data: "new-data",
					})
					require.NoError(t, err)
					require.True(t, ok)

					vectors, err := client.Fetch(Fetch{
						Ids:         []string{id},
						IncludeData: true,
					})
					require.NoError(t, err)
					require.Equal(t, 1, len(vectors))
					require.Equal(t, id, vectors[0].Id)
					require.Equal(t, "new-data", vectors[0].Data)
				})

				t.Run("metadata", func(t *testing.T) {
					id := randomString()
					v, sv := randomVectors(tcType)

					err := client.Upsert(Upsert{
						Id:           id,
						Vector:       v,
						SparseVector: sv,
						Metadata:     map[string]any{"old": "metadata"},
					})
					require.NoError(t, err)

					ok, err := client.Update(Update{
						Id:       id,
						Metadata: map[string]any{"new": "metadata"},
					})
					require.NoError(t, err)
					require.True(t, ok)

					vectors, err := client.Fetch(Fetch{
						Ids:             []string{id},
						IncludeMetadata: true,
					})
					require.NoError(t, err)
					require.Equal(t, 1, len(vectors))
					require.Equal(t, id, vectors[0].Id)
					require.Equal(t, map[string]any{"new": "metadata"}, vectors[0].Metadata)
				})

				t.Run("patch metadata", func(t *testing.T) {
					id := randomString()
					v, sv := randomVectors(tcType)

					err := client.Upsert(Upsert{
						Id:           id,
						Vector:       v,
						SparseVector: sv,
						Metadata:     map[string]any{"old": "metadata"},
					})
					require.NoError(t, err)

					ok, err := client.Update(Update{
						Id:                 id,
						Metadata:           map[string]any{"new": "metadata"},
						MetadataUpdateMode: MetadataUpdateModePatch,
					})
					require.NoError(t, err)
					require.True(t, ok)

					vectors, err := client.Fetch(Fetch{
						Ids:             []string{id},
						IncludeMetadata: true,
					})
					require.NoError(t, err)
					require.Equal(t, 1, len(vectors))
					require.Equal(t, id, vectors[0].Id)
					require.Equal(t, map[string]any{"old": "metadata", "new": "metadata"}, vectors[0].Metadata)
				})

				t.Run("vector", func(t *testing.T) {
					id := randomString()
					v0, sv0 := randomVectors(tcType)

					err := client.Upsert(Upsert{
						Id:           id,
						Vector:       v0,
						SparseVector: sv0,
					})
					require.NoError(t, err)

					v1, sv1 := randomVectors(tcType)

					ok, err := client.Update(Update{
						Id:           id,
						Vector:       v1,
						SparseVector: sv1,
					})
					require.NoError(t, err)
					require.True(t, ok)

					vectors, err := client.Fetch(Fetch{
						Ids:            []string{id},
						IncludeVectors: true,
					})
					require.NoError(t, err)
					require.Equal(t, 1, len(vectors))
					require.Equal(t, id, vectors[0].Id)
					require.Equal(t, v1, vectors[0].Vector)
					require.Equal(t, sv1, vectors[0].SparseVector)
				})
			})
		}
	}
}
