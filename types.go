package vector

type Upsert struct {
	Id       string         `json:"id"`
	Vector   []float32      `json:"vector"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Query struct {
	Vector          []float32 `json:"vector"`
	TopK            int       `json:"topK,omitempty"`
	IncludeVectors  bool      `json:"includeVectors,omitempty"`
	IncludeMetadata bool      `json:"includeMetadata,omitempty"`
}

type Vector struct {
	Id       string         `json:"id"`
	Vector   []float32      `json:"vector,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type VectorScore struct {
	Id       string         `json:"id"`
	Score    float32        `json:"score"`
	Vector   []float32      `json:"vector,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Fetch struct {
	Ids             []string `json:"ids"`
	IncludeVectors  bool     `json:"includeVectors,omitempty"`
	IncludeMetadata bool     `json:"includeMetadata,omitempty"`
}

type Range struct {
	Cursor          string `json:"cursor,omitempty"`
	Limit           int    `json:"limit,omitempty"`
	IncludeVectors  bool   `json:"includeVectors,omitempty"`
	IncludeMetadata bool   `json:"includeMetadata,omitempty"`
}

type RangeVectors struct {
	NextCursor string   `json:"nextCursor"`
	Vectors    []Vector `json:"vectors,omitempty"`
}

type deleted struct {
	Deleted int `json:"deleted"`
}

type IndexInfo struct {
	VectorCount        int    `json:"vectorCount"`
	PendingVectorCount int    `json:"pendingVectorCount"`
	IndexSize          int    `json:"indexSize"`
	Dimension          int    `json:"dimension"`
	SimilarityFunction string `json:"similarityFunction"`
}

type Response[T any] struct {
	Result T      `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
}
