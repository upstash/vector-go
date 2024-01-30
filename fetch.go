package vector

const fetchPath = "/fetch"

// Fetch fetches one or more vectors with the ids passed into f.
// When f.IncludeVectors is true, values of the vectors are also
// returned. When f.IncludeMetadata is true, metadata of the vectors
// are also returned, if any.
func (ix *Index) Fetch(f Fetch) (vectors []Vector, err error) {
	data, err := ix.sendJson(fetchPath, f)
	if err != nil {
		return
	}
	vectors, err = parseResponse[[]Vector](data)
	return
}
