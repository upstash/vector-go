package vector

const queryPath = "/query"

// Query returns the result of the query for the given vector in the default namespace.
// When q.TopK is specified, the result will contain at most q.TopK many vectors.
// The returned list will contain vectors sorted in descending order of score,
// which correlates with the similarity of the vectors to the given query vector.
// When q.IncludeVectors is true, values of the vectors are also returned.
// When q.IncludeMetadata is true, metadata of the vectors are also returned, if any.
func (ix *Index) Query(q Query) (scores []VectorScore, err error) {
	return ix.queryInternal(q, defaultNamespace)
}

func (ix *Index) queryInternal(q Query, ns string) (scores []VectorScore, err error) {
	data, err := ix.sendJson(buildPath(queryPath, ns), q)
	if err != nil {
		return
	}
	scores, err = parseResponse[[]VectorScore](data)
	return
}
