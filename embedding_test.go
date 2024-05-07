package vector

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEmbedding(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newEmbeddingTestClient(ns)
			require.NoError(t, err)

			id0 := "tr"
			id1 := "jp"
			id2 := "uk"
			id3 := "fr"
			err = client.UpsertDataMany([]UpsertData{
				{
					Id:       id0,
					Data:     "Capital of Türkiye is Ankara.",
					Metadata: map[string]any{"country": id0, "capital": "Ankara"},
				},
				{
					Id:       id1,
					Data:     "Capital of Japan is Tokyo.",
					Metadata: map[string]any{"country": id1, "capital": "Tokyo"},
				},
				{
					Id:       id2,
					Data:     "Capital of England is London.",
					Metadata: map[string]any{"country": id2, "capital": "London"},
				},
				{
					Id:       id3,
					Data:     "Capital of France is Paris.",
					Metadata: map[string]any{"country": id3, "capital": "Paris"},
				},
			})
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				info, err := client.Info()
				require.NoError(t, err)
				return info.PendingVectorCount == 0
			}, 10*time.Second, 1*time.Second)

			t.Run("score", func(t *testing.T) {
				scores, err := client.QueryData(QueryData{
					Data: "where is the capital of Japan?",
					TopK: 1,
				})
				require.NoError(t, err)
				require.Equal(t, 1, len(scores))
				require.Equal(t, id1, scores[0].Id)
			})

			t.Run("with metadata", func(t *testing.T) {
				scores, err := client.QueryData(QueryData{
					Data:            "Which country's capital is Ankara?",
					TopK:            1,
					IncludeMetadata: true,
				})
				require.NoError(t, err)
				require.Equal(t, 1, len(scores))
				require.Equal(t, id0, scores[0].Id)
				require.Equal(t, map[string]any{"country": "tr", "capital": "Ankara"}, scores[0].Metadata)
			})

			t.Run("with metadata filtering", func(t *testing.T) {
				query := QueryData{
					Data:            "Where is the capital of France?",
					TopK:            1,
					IncludeMetadata: true,
					Filter:          `country = 'fr'`,
				}

				scores, err := client.QueryData(query)
				require.NoError(t, err)
				require.Equal(t, 1, len(scores))
				require.Equal(t, id3, scores[0].Id)
				require.Equal(t, map[string]any{"country": "fr", "capital": "Paris"}, scores[0].Metadata)
			})
		})
	}
}
