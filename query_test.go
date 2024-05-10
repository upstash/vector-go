package vector

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient()
			require.NoError(t, err)

			namespace := client.Namespace(ns)
			id0 := randomString()
			id1 := randomString()
			id2 := randomString()
			err = namespace.UpsertMany([]Upsert{
				{
					Id:       id0,
					Vector:   []float32{0, 1},
					Metadata: map[string]any{"foo": "bar"},
				},
				{
					Id:     id1,
					Vector: []float32{5, 10},
				},
				{
					Id:       id2,
					Vector:   []float32{0.01, 1.01},
					Metadata: map[string]any{"foo": "nay"},
				},
			})
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				info, err := client.Info()
				require.NoError(t, err)
				return info.PendingVectorCount == 0
			}, 10*time.Second, 1*time.Second)

			t.Run("score", func(t *testing.T) {
				scores, err := namespace.Query(Query{
					Vector: []float32{0, 1},
					TopK:   2,
				})
				require.NoError(t, err)
				require.Equal(t, 2, len(scores))
				require.Equal(t, id0, scores[0].Id)
				require.Equal(t, float32(1.0), scores[0].Score)
				require.Equal(t, id2, scores[1].Id)
			})

			t.Run("with metadata and vectors", func(t *testing.T) {
				scores, err := namespace.Query(Query{
					Vector:          []float32{0, 1},
					TopK:            2,
					IncludeMetadata: true,
					IncludeVectors:  true,
				})
				require.NoError(t, err)
				require.Equal(t, 2, len(scores))
				require.Equal(t, id0, scores[0].Id)
				require.Equal(t, float32(1.0), scores[0].Score)
				require.Equal(t, map[string]any{"foo": "bar"}, scores[0].Metadata)
				require.Equal(t, []float32{0, 1}, scores[0].Vector)

				require.Equal(t, id2, scores[1].Id)
				require.Equal(t, []float32{0.01, 1.01}, scores[1].Vector)
			})

			t.Run("with metadata filtering", func(t *testing.T) {
				query := Query{
					Vector:          []float32{0, 1},
					TopK:            10,
					IncludeMetadata: true,
					IncludeVectors:  true,
					Filter:          `foo = 'bar'`,
				}

				scores, err := namespace.Query(query)
				require.NoError(t, err)
				require.Equal(t, 1, len(scores))
				require.Equal(t, id0, scores[0].Id)
				require.Equal(t, float32(1.0), scores[0].Score)
				require.Equal(t, map[string]any{"foo": "bar"}, scores[0].Metadata)
				require.Equal(t, []float32{0, 1}, scores[0].Vector)

				query.Filter = `foo = 'nay'`
				scores, err = namespace.Query(query)
				require.NoError(t, err)
				require.Equal(t, 1, len(scores))
				require.Equal(t, id2, scores[0].Id)
				require.Equal(t, map[string]any{"foo": "nay"}, scores[0].Metadata)
				require.Equal(t, []float32{0.01, 1.01}, scores[0].Vector)
			})
		})
	}
}
