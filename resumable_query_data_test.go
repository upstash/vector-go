package vector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResumableQueryData(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeDenseEmbedding, ns)
			require.NoError(t, err)

			u := make([]UpsertData, 10)
			for i := 0; i < 10; i++ {
				id := randomString()
				u[i] = UpsertData{
					Id:       id,
					Data:     "Upstash",
					Metadata: map[string]any{"metadata": id},
				}
			}

			err = client.UpsertDataMany(u)
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				info, err := client.Info()
				require.NoError(t, err)
				return info.PendingVectorCount == 0
			}, 10*time.Second, 1*time.Second)

			t.Run("next", func(t *testing.T) {
				scores, handle, err := client.ResumableQueryData(ResumableQueryData{
					Data: "Upstash",
					TopK: 2,
				})

				t.Cleanup(func() {
					if handle != nil {
						handle.Close()
					}
				})

				require.NoError(t, err)
				require.Equal(t, 2, len(scores))

				scores, err = handle.Next(ResumableQueryNext{
					AdditionalK: 3,
				})
				require.NoError(t, err)
				require.Equal(t, 3, len(scores))

				scores, err = handle.Next(ResumableQueryNext{
					AdditionalK: 4,
				})
				require.NoError(t, err)
				require.Equal(t, 4, len(scores))

				err = handle.Close()
				require.NoError(t, err)
			})

			t.Run("next with metadata and vectors", func(t *testing.T) {
				validateScores := func(scores []VectorScore) {
					for _, score := range scores {
						id := score.Id
						require.Equal(t, map[string]any{"metadata": id}, score.Metadata)
						require.Equal(t, "Upstash", score.Data)
					}
				}

				scores, handle, err := client.ResumableQueryData(ResumableQueryData{
					Data:            "Upstash",
					TopK:            2,
					IncludeMetadata: true,
					IncludeData:     true,
				})

				t.Cleanup(func() {
					if handle != nil {
						handle.Close()
					}
				})

				require.NoError(t, err)
				require.Equal(t, 2, len(scores))
				validateScores(scores)

				scores, err = handle.Next(ResumableQueryNext{
					AdditionalK: 3,
				})
				require.NoError(t, err)
				require.Equal(t, 3, len(scores))
				validateScores(scores)

				scores, err = handle.Next(ResumableQueryNext{
					AdditionalK: 4,
				})
				require.NoError(t, err)
				require.Equal(t, 4, len(scores))
				validateScores(scores)

				err = handle.Close()
				require.NoError(t, err)
			})
		})
	}
}
