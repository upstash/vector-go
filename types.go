package vector

type Upsert struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Vector values.
	Vector []float32 `json:"vector"`

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

type Query struct {
	// The query vector.
	Vector []float32 `json:"vector"`

	// The maximum number of vectors that will
	// be returned for the query response.
	TopK int `json:"topK,omitempty"`

	// Whether to include vector values in the query response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the query response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`

	// Query filter
	Filter any `json:"filter,omitempty"`
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

	// Query filter
	Filter any `json:"filter,omitempty"`
}

type Vector struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Vector values.
	Vector []float32 `json:"vector,omitempty"`

	// Optional metadata of the vector, if any.
	Metadata map[string]any `json:"metadata,omitempty"`
}

type VectorScore struct {
	// Unique id of the vector.
	Id string `json:"id"`

	// Similarity score of the vector to the query vector.
	// Vectors more similar to the query vector have higher score.
	Score float32 `json:"score"`

	// Optional vector values.
	Vector []float32 `json:"vector,omitempty"`

	// Optional metadata of the vector, if any.
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Fetch struct {
	// Unique vectors ids to fetch.
	Ids []string `json:"ids"`

	// Whether to include vector values in the fetch response.
	IncludeVectors bool `json:"includeVectors,omitempty"`

	// Whether to include metadata in the fetch response, if any.
	IncludeMetadata bool `json:"includeMetadata,omitempty"`
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

	NamespaceInfo map[string]NamespaceInfo `json:"namespaces"`
}

type NamespaceInfo struct {
	// The number of vectors in the index.
	VectorCount int `json:"vectorCount"`

	// The number of vectors that are pending to be indexed.
	PendingVectorCount int `json:"pendingVectorCount"`
}

type response[T any] struct {
	Result T      `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
}
