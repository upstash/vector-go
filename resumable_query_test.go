package vector

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResumableQuery(t *testing.T) {
	makeSparseVectorsConnected := func(sv *SparseVector) {
		if sv == nil {
			return
		}

		// add at least one common non-zero dimension
		// to make sure that the tests can find
		// results to return for sparse vectors.
		if !slices.Contains(sv.Indices, 0) {
			sv.Indices = append(sv.Indices, 0)
			sv.Values = append(sv.Values, 0.5)
		}
	}

	for _, ns := range namespaces {
		for _, tcType := range testClientTypes {
			t.Run("namespace_"+ns+"_index_type_"+string(tcType), func(t *testing.T) {
				client, err := newTestClient(tcType, ns)
				require.NoError(t, err)

				u := make([]Upsert, 10)
				for i := 0; i < 10; i++ {
					v, sv := randomVectors(tcType)
					makeSparseVectorsConnected(sv)

					id := randomString()
					u[i] = Upsert{
						Id:           id,
						Vector:       v,
						SparseVector: sv,
						Metadata:     map[string]any{"metadata": id},
						Data:         id + "-data",
					}
				}

				err = client.UpsertMany(u)
				require.NoError(t, err)

				require.Eventually(t, func() bool {
					info, err := client.Info()
					require.NoError(t, err)
					return info.PendingVectorCount == 0
				}, 10*time.Second, 1*time.Second)

				t.Run("next", func(t *testing.T) {
					v, sv := randomVectors(tcType)
					makeSparseVectorsConnected(sv)

					scores, handle, err := client.ResumableQuery(ResumableQuery{
						Vector:       v,
						SparseVector: sv,
						TopK:         2,
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
					v, sv := randomVectors(tcType)
					makeSparseVectorsConnected(sv)

					validateScores := func(scores []VectorScore) {
						for _, score := range scores {
							id := score.Id
							require.Equal(t, map[string]any{"metadata": id}, score.Metadata)
							require.Equal(t, id+"-data", score.Data)
							switch tcType {
							case testClientTypeDense:
								require.NotNil(t, score.Vector)
							case testClientTypeSparse:
								require.NotNil(t, score.SparseVector)
							default:
								require.NotNil(t, score.Vector)
								require.NotNil(t, score.SparseVector)
							}
						}
					}

					scores, handle, err := client.ResumableQuery(ResumableQuery{
						Vector:          v,
						SparseVector:    sv,
						TopK:            2,
						IncludeVectors:  true,
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
}
