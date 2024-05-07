package vector

const queryDataPath = "/query-data"

// QueryData returns the result of the query for the given data
// by converting it to an embedding on the server.
// When q.TopK is specified, the result will contain at most q.TopK
// many vectors. The returned list will contain vectors sorted in descending
// order of score, which correlates with the similarity of the vectors to the
// given query vector. When q.IncludeVectors is true, values of the vectors are
// also returned. When q.IncludeMetadata is true, metadata of the vectors are
// also returned, if any.
func (ix *Index) QueryData(q QueryData) (scores []VectorScore, err error) {
	data, err := ix.sendJson(queryDataPath, q, true)
	if err != nil {
		return
	}
	scores, err = parseResponse[[]VectorScore](data)
	return
}
