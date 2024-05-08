package vector

const fetchPath = "/fetch"

// Fetch fetches one or more vectors in a namespace with the ids passed into f.
// If IncludeVectors is set to true, the vector values are also returned.
// If IncludeMetadata is set to true, any associated metadata of the vectors is also returned, if any.
// If namespace is not specified, the default namespace is used.
func (ix *Index) Fetch(f Fetch) (vectors []Vector, err error) {
	data, err := ix.sendJson(fetchPath, f, true)
	if err != nil {
		return
	}
	vectors, err = parseResponse[[]Vector](data)
	return
}
