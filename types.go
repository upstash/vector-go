package vector

type SparseVector struct {
	// List of dimensions that have non-zero values.
	Indices []int32 `json:"indices"`

	// Values of the non-zero dimensions.
	// It must be of the same size as the indices.
	Values []float32 `json:"values"`
}

type Upsert struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Dense vector values for dense and hybrid indexes.
	Vector []float32 `json:"vector,omitempty"`

	// Sparse vector values for sparse and hybrid indexes.
	SparseVector *SparseVector `json:"sparseVector,omitempty"`

	// Optional data of the vector.
	Data string `json:"data,omitempty"`

	// Optional metadata of the vector.
	Metadata map[string]any `json:"metadata,omitempty"`
}

type UpsertData struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Raw data.
	// Data will be converted to the vector embedding on the server.
	Data string `json:"data"`

	// Optional metadata of the vector.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// WeightingStrategy specifies what kind of weighting strategy
// should be used while querying the matching non-zero
// dimension values of the query vector with the indexed vectors
// for sparse vectors.
type WeightingStrategy string

const (
	// WeightingStrategyIDF uses inverse document frequencies.
	//
	// It is recommended to use this weighting strategy for
	// BM25 sparse embedding models.
	//
	// It is calculated as
	// ln(((N - n(q) + 0.5) / (n(q) + 0.5)) + 1) where
	// N:    Total number of sparse vectors.
	// n(q): Total number of sparse vectors having non-zero value
	// for that particular dimension.
	// ln:   Natural logarithm
	// The values of N and n(q) are maintained by Upstash as the
	// vectors are indexed.
	WeightingStrategyIDF WeightingStrategy = "IDF"
)

// FusionAlgorithm specifies the algorithm to use while fusing scores
// from dense and sparse components of a hybrid index.
type FusionAlgorithm string

const (
	// FusionAlgorithmRRF is reciprocal rank fusion.
	//
	// Each sorted score from the dense and sparse indexes are
	// mapped to 1 / (rank + K), where rank is the order of the
	// score in the dense or sparse scores and K is a constant
	// with the value of 60.
	//
	// Then, scores from the dense and sparse components are
	// deduplicated (i.e. if a score for the same vector is present
	// in both dense and sparse scores, the mapped scores are
	// added; otherwise individual mapped scores are used)
	// and the final result is returned as the topK values
	// of this final list.
	//
	// In short, this algorithm just takes the order of the scores
	// into consideration.
	FusionAlgorithmRRF FusionAlgorithm = "RRF"

	// FusionAlgorithmDBSF is distribution based score fusion.
	//
	// Each sorted score from the dense and sparse indexes are
	// normalized as
	// (s - (mean - 3 * stddev)) / ((mean + 3 * stddev) - (mean - 3 * stddev))
	// where s is the score, (mean - 3 * stddev) is the minimum,
	// and (mean + 3 * stddev) is the maximum tail ends of the distribution.
	//
	// Then, scores from the dense and sparse components are
	// deduplicated (i.e. if a score for the same vector is present
	// in both dense and sparse scores, the normalized scores are
	// added; otherwise individual normalized scores are used)
	// and the final result is returned as the topK values
	// of this final list.
	//
	// In short, this algorithm takes distribution of the scores
	// into consideration as well, as opposed to the RRF.
	FusionAlgorithmDBSF FusionAlgorithm = "DBSF"
)

// QueryMode for hybrid indexes with Upstash-hosted embedding models.
//
// It specifies whether to run the query in only the
// dense index, only the sparse index, or in both.
type QueryMode string

const (
	// QueryModeHybrid runs the query in hybrid index mode,
	// after embedding the raw text data into dense and sparse vectors.
	//
	// Query results from the dense and sparse index components
	// of the hybrid index are fused before returning the result.
	QueryModeHybrid QueryMode = "HYBRID"

	// QueryModeDense runs the query in dense index mode,
	// after embedding the raw text data into a dense vector.
	//
	// Only the query results from the dense index component
	// of the hybrid index is returned.
	QueryModeDense QueryMode = "DENSE"

	// QueryModeSparse runs the query in sparse index mode,
	// after embedding the raw text data into a sparse vector.
	//
	// Only the query results from the sparse index component
	// of the hybrid index is returned.
	QueryModeSparse QueryMode = "SPARSE"
)

type Query struct {
	// The dense query vector for dense and hybrid indexes.
	Vector []float32 `json:"vector,omitempty"`

	// The sparse query vector for sparse and hybrid indexes.
	SparseVector *SparseVector `json:"sparseVector,omitempty"`

	// The maximum number of vectors that will
	// be returned for the query response.
	TopK int `json:"topK,omitempty"`

	// Whether to include vector values in the query response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the query response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Whether to include data in the query response, if any.
	IncludeData bool `json:"includeData,omitempty"`

	// Query filter
	Filter any `json:"filter,omitempty"`

	// Weighting strategy to be used for sparse vectors.
	// If not provided, no weighting will be used.
	WeightingStrategy WeightingStrategy `json:"weightingStrategy,omitempty"`

	// Fusion algorithm to use while fusing scores
	// from dense and sparse components of a hybrid index.
	// If not provided, defaults to RRF.
	FusionAlgorithm FusionAlgorithm `json:"fusionAlgorithm,omitempty"`
}

type QueryData struct {
	// Raw data.
	// Data will be converted to the vector embedding on the server.
	Data string `json:"data"`

	// The maximum number of vectors that will
	// be returned for the query response.
	TopK int `json:"topK,omitempty"`

	// Whether to include vector values in the query response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the query response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Whether to include data in the query response, if any.
	IncludeData bool `json:"includeData,omitempty"`

	// Query filter
	Filter any `json:"filter,omitempty"`

	// Weighting strategy to be used for sparse vectors.
	// If not provided, no weighting will be used.
	WeightingStrategy WeightingStrategy `json:"weightingStrategy,omitempty"`

	// Fusion algorithm to use while fusing scores
	// from dense and sparse components of a hybrid index.
	// If not provided, defaults to RRF.
	FusionAlgorithm FusionAlgorithm `json:"fusionAlgorithm,omitempty"`

	// Specifies whether to run the query in only the
	// dense index, only the sparse index, or in both for hybrid
	// indexes with Upstash-hosted embedding models.
	// If not provided, defaults to hybrid query mode.
	QueryMode QueryMode `json:"queryMode,omitempty"`
}

type Vector struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Dense vector values for dense and hybrid indexes.
	Vector []float32 `json:"vector,omitempty"`

	// Sparse vector values for sparse and hybrid indexes.
	SparseVector *SparseVector `json:"sparseVector,omitempty"`

	// Optional metadata of the vector, if any.
	Metadata map[string]any `json:"metadata,omitempty"`

	// Optional data of the vector.
	Data string `json:"data,omitempty"`
}

type VectorScore struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Similarity score of the vector to the query vector.
	// Vectors more similar to the query vector have higher score.
	Score float32 `json:"score"`

	// Optional dense vector values for dense and hybrid indexes.
	Vector []float32 `json:"vector,omitempty"`

	// Optional sparse vector values for sparse and hybrid indexes.
	SparseVector *SparseVector `json:"sparseVector,omitempty"`

	// Optional metadata of the vector, if any.
	Metadata map[string]any `json:"metadata,omitempty"`

	// Optional data of the vector.
	Data string `json:"data,omitempty"`
}

type Fetch struct {
	// Unique vectors ids to fetch.
	Ids []string `json:"ids"`

	// Whether to include vector values in the fetch response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the fetch response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Whether to include data in the query response, if any.
	IncludeData bool `json:"includeData,omitempty"`
}

type Range struct {
	// The cursor to start returning range from (inclusive).
	Cursor string `json:"cursor,omitempty"`

	// The maximum number of vectors that will be returned for
	// the range response.
	Limit int `json:"limit,omitempty"`

	// Whether to include vector values in the range response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the fetch response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Whether to include data in the query response, if any.
	IncludeData bool `json:"includeData,omitempty"`
}

type RangeVectors struct {
	// The cursor that should be used for the subsequent range requests.
	NextCursor string `json:"nextCursor"`

	// List of vectors in the range.
	Vectors []Vector `json:"vectors,omitempty"`
}

type deleted struct {
	Deleted int `json:"deleted"`
}

type IndexInfo struct {
	// The number of vectors in the index.
	VectorCount int `json:"vectorCount"`

	// The number of vectors that are pending to be indexed.
	PendingVectorCount int `json:"pendingVectorCount"`

	// The size of the index on disk in bytes
	IndexSize int `json:"indexSize"`

	// The dimension of the vectors.
	Dimension int `json:"dimension"`

	// Name of the similarity function used in indexing and queries.
	SimilarityFunction string `json:"similarityFunction"`

	// Per-namespace vector and pending vector counts
	Namespaces map[string]NamespaceInfo `json:"namespaces"`
}

type NamespaceInfo struct {
	// The number of vectors in the namespace of the index.
	VectorCount int `json:"vectorCount"`

	// The number of vectors that are pending to be indexed.
	PendingVectorCount int `json:"pendingVectorCount"`
}

// MetadataUpdateMode specifies whether to overwrite the whole
// metadata while updating it, or patch the metadata
// (insert new fields or update or delete existing fields)
// according to the RFC 7396 JSON Merge Patch algorithm.
type MetadataUpdateMode string

const (
	// MetadataUpdateModeOverwrite overwrites the metadata,
	// and set it to a new value.
	MetadataUpdateModeOverwrite MetadataUpdateMode = "OVERWRITE"

	// MetadataUpdateModePatch patches the metadata according
	// to the JSON Merge Patch algorithm.
	MetadataUpdateModePatch MetadataUpdateMode = "PATCH"
)

type Update struct {
	// The id of the vector to update.
	Id string `json:"id"`

	// The new dense vector values for dense and hybrid indexes.
	Vector []float32 `json:"vector,omitempty"`

	// The new sparse vector values for sparse and hybrid indexes.
	SparseVector *SparseVector `json:"sparseVector,omitempty"`

	// The new data of the vector.
	Data string `json:"data,omitempty"`

	// The new metadata of the vector.
	Metadata map[string]any `json:"metadata,omitempty"`

	// Whether to overwrite the whole metadata while updating it,
	// or patch the metadata according to the JSON Merge Patch algorithm.
	// If not provided, defaults to overwrite.
	MetadataUpdateMode MetadataUpdateMode `json:"metadataUpdateMode,omitempty"`
}

type updated struct {
	Updated int `json:"updated"`
}

type response[T any] struct {
	Result T      `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
}

type ResumableQuery struct {
	// The dense query vector for dense and hybrid indexes.
	Vector []float32 `json:"vector,omitempty"`

	// The sparse query vector for sparse and hybrid indexes.
	SparseVector *SparseVector `json:"sparseVector,omitempty"`

	// The maximum number of vectors that will
	// be returned for the query response.
	TopK int `json:"topK,omitempty"`

	// Whether to include vector values in the query response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the query response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Whether to include data in the query response, if any.
	IncludeData bool `json:"includeData,omitempty"`

	// Query filter
	Filter any `json:"filter,omitempty"`

	// Weighting strategy to be used for sparse vectors.
	// If not provided, no weighting will be used.
	WeightingStrategy WeightingStrategy `json:"weightingStrategy,omitempty"`

	// Fusion algorithm to use while fusing scores
	// from dense and sparse components of a hybrid index.
	// If not provided, defaults to RRF.
	FusionAlgorithm FusionAlgorithm `json:"fusionAlgorithm,omitempty"`

	// Maximum idle time for the resumable query in seconds.
	// If not provided, defaults to 1 hour.
	MaxIdle uint32 `json:"maxIdle,omitempty"`
}

type ResumableQueryData struct {
	// Raw data.
	// Data will be converted to the vector embedding on the server.
	Data string `json:"data"`

	// The maximum number of vectors that will
	// be returned for the query response.
	TopK int `json:"topK,omitempty"`

	// Whether to include vector values in the query response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the query response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Whether to include data in the query response, if any.
	IncludeData bool `json:"includeData,omitempty"`

	// Query filter
	Filter any `json:"filter,omitempty"`

	// Weighting strategy to be used for sparse vectors.
	// If not provided, no weighting will be used.
	WeightingStrategy WeightingStrategy `json:"weightingStrategy,omitempty"`

	// Fusion algorithm to use while fusing scores
	// from dense and sparse components of a hybrid index.
	// If not provided, defaults to RRF.
	FusionAlgorithm FusionAlgorithm `json:"fusionAlgorithm,omitempty"`

	// Specifies whether to run the query in only the
	// dense index, only the sparse index, or in both for hybrid
	// indexes with Upstash-hosted embedding models.
	// If not provided, defaults to hybrid query mode.
	QueryMode QueryMode `json:"queryMode,omitempty"`

	// Maximum idle time for the resumable query in seconds.
	// If not provided, defaults to 1 hour.
	MaxIdle uint32 `json:"maxIdle,omitempty"`
}

type ResumableQueryNext struct {
	AdditionalK int `json:"additionalK,omitempty"`
}

type resumableQueryStart struct {
	UUID   string        `json:"uuid"`
	Scores []VectorScore `json:"scores"`
}

type resumableQueryNext struct {
	ResumableQueryNext
	UUID string `json:"uuid"`
}

type resumableQueryEnd struct {
	UUID string `json:"uuid"`
}
