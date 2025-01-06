package vector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	for _, ns := range namespaces {
		for _, tcType := range testClientTypes {
			t.Run("namespace_"+ns+"_index_type_"+string(tcType), func(t *testing.T) {
				client, err := newTestClient(tcType, ns)
				require.NoError(t, err)

				id0 := randomString()
				v0, sv0 := randomVectors(tcType)

				id1 := randomString()
				// make v1 and sv1 similar to v0 and sv0 so that
				// they would be selected as the second most similar
				// vector when we use v0 and sv0 for queries.
				var v1 []float32
				var sv1 *SparseVector

				if v0 != nil {
					v1 = make([]float32, len(v0))
					copy(v1, v0)
					for i, val := range v1 {
						v1[i] = val + 0.01
					}
				}

				if sv0 != nil {
					sv1Indices := make([]int32, len(sv0.Indices))
					copy(sv1Indices, sv0.Indices)

					sv1Values := make([]float32, len(sv0.Values))
					copy(sv1Values, sv0.Values)
					for i, val := range sv1Values {
						sv1Values[i] = max(0.0001, val-0.1)
					}

					sv1 = &SparseVector{
						Indices: sv1Indices,
						Values:  sv1Values,
					}
				}

				id2 := randomString()
				v2, sv2 := randomVectors(tcType)

				err = client.UpsertMany([]Upsert{
					{
						Id:           id0,
						Vector:       v0,
						SparseVector: sv0,
						Metadata:     map[string]any{"foo": "bar"},
						Data:         "vector0 data",
					},
					{
						Id:           id1,
						Vector:       v1,
						SparseVector: sv1,
						Data:         "vector1 data",
					},
					{
						Id:           id2,
						Vector:       v2,
						SparseVector: sv2,
						Metadata:     map[string]any{"foo": "nay"},
					},
				})
				require.NoError(t, err)

				require.Eventually(t, func() bool {
					info, err := client.Info()
					require.NoError(t, err)
					return info.PendingVectorCount == 0
				}, 10*time.Second, 1*time.Second)

				t.Run("score", func(t *testing.T) {
					scores, err := client.Query(Query{
						Vector:       v0,
						SparseVector: sv0,
						TopK:         2,
					})
					require.NoError(t, err)
					require.Equal(t, 2, len(scores))
					require.Equal(t, id0, scores[0].Id)
					if tcType == testClientTypeDense {
						require.Equal(t, float32(1.0), scores[0].Score)
					}
					require.Equal(t, id1, scores[1].Id)
				})

				t.Run("with metadata and vectors", func(t *testing.T) {
					scores, err := client.Query(Query{
						Vector:          v0,
						SparseVector:    sv0,
						TopK:            2,
						IncludeMetadata: true,
						IncludeVectors:  true,
						IncludeData:     true,
					})
					require.NoError(t, err)
					require.Equal(t, 2, len(scores))
					require.Equal(t, id0, scores[0].Id)
					if tcType == testClientTypeDense {
						require.Equal(t, float32(1.0), scores[0].Score)
					}
					require.Equal(t, map[string]any{"foo": "bar"}, scores[0].Metadata)
					require.Equal(t, "vector0 data", scores[0].Data)
					require.Equal(t, v0, scores[0].Vector)
					require.Equal(t, sv0, scores[0].SparseVector)

					require.Equal(t, id1, scores[1].Id)
					require.Equal(t, v1, scores[1].Vector)
					require.Equal(t, sv1, scores[1].SparseVector)
					require.Empty(t, scores[1].Metadata)
					require.Equal(t, "vector1 data", scores[1].Data)
				})

				t.Run("with metadata filtering", func(t *testing.T) {
					query := Query{
						Vector:          v0,
						SparseVector:    sv0,
						TopK:            10,
						IncludeMetadata: true,
						IncludeVectors:  true,
						IncludeData:     true,
						Filter:          `foo = 'bar'`,
					}

					scores, err := client.Query(query)
					require.NoError(t, err)
					require.Equal(t, 1, len(scores))
					require.Equal(t, id0, scores[0].Id)
					if tcType == testClientTypeDense {
						require.Equal(t, float32(1.0), scores[0].Score)
					}
					require.Equal(t, map[string]any{"foo": "bar"}, scores[0].Metadata)
					require.Equal(t, "vector0 data", scores[0].Data)
					require.Equal(t, v0, scores[0].Vector)
					require.Equal(t, sv0, scores[0].SparseVector)

					query.Vector = v2
					query.SparseVector = sv2
					query.Filter = `foo = 'nay'`
					scores, err = client.Query(query)
					require.NoError(t, err)
					require.Equal(t, 1, len(scores))
					require.Equal(t, id2, scores[0].Id)
					require.Equal(t, map[string]any{"foo": "nay"}, scores[0].Metadata)
					require.Equal(t, v2, scores[0].Vector)
					require.Equal(t, sv2, scores[0].SparseVector)
					require.Empty(t, scores[0].Data)
				})
			})
		}
	}
}

func TestQueryWeightingStrategy(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeSparse, ns)
			require.NoError(t, err)

			err = client.UpsertMany([]Upsert{
				{
					Id: "id0",
					SparseVector: &SparseVector{
						Indices: []int32{0, 1},
						Values:  []float32{0.1, 0.1},
					},
				},
				{
					Id: "id1",
					SparseVector: &SparseVector{
						Indices: []int32{1, 2},
						Values:  []float32{0.1, 0.1},
					},
				},
				{
					Id: "id2",
					SparseVector: &SparseVector{
						Indices: []int32{2, 3},
						Values:  []float32{0.1, 0.1},
					},
				},
			})
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				info, err := client.Info()
				require.NoError(t, err)
				return info.PendingVectorCount == 0
			}, 10*time.Second, 1*time.Second)

			scores, err := client.Query(Query{
				SparseVector: &SparseVector{
					Indices: []int32{0, 1, 3},
					Values:  []float32{0.2, 0.1, 0.1},
				},
				WeightingStrategy: WeightingStrategyIDF,
			})
			require.NoError(t, err)

			require.Equal(t, 3, len(scores))
			require.Equal(t, "id0", scores[0].Id)
			require.Equal(t, "id2", scores[1].Id)
			require.Equal(t, "id1", scores[2].Id)
		})
	}
}

func TestQueryFusionAlgorithm(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeHybrid, ns)
			require.NoError(t, err)

			err = client.UpsertMany([]Upsert{
				{
					Id:     "id0",
					Vector: []float32{0.8, 0.9},
					SparseVector: &SparseVector{
						Indices: []int32{0, 1},
						Values:  []float32{0.1, 0.1},
					},
				},
				{
					Id:     "id1",
					Vector: []float32{0.9, 0.9},
					SparseVector: &SparseVector{
						Indices: []int32{1, 2},
						Values:  []float32{0.1, 0.1},
					},
				},
				{
					Id:     "id2",
					Vector: []float32{0.3, 0.9},
					SparseVector: &SparseVector{
						Indices: []int32{2, 3},
						Values:  []float32{0.1, 0.1},
					},
				},
			})
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				info, err := client.Info()
				require.NoError(t, err)
				return info.PendingVectorCount == 0
			}, 10*time.Second, 1*time.Second)

			scores, err := client.Query(Query{
				Vector: []float32{0.9, 0.9},
				SparseVector: &SparseVector{
					Indices: []int32{0, 1, 3},
					Values:  []float32{0.1, 0.1, 100},
				},
				FusionAlgorithm: FusionAlgorithmDBSF,
			})
			require.NoError(t, err)

			require.Equal(t, 3, len(scores))
			require.Equal(t, "id1", scores[0].Id)
			require.Equal(t, "id2", scores[1].Id)
			require.Equal(t, "id0", scores[2].Id)
		})
	}
}
