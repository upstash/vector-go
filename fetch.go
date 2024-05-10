package vector

const fetchPath = "/fetch"

// Fetch fetches one or more vectors in the default namespace with the ids passed into f.
// If IncludeVectors is set to true, the vector values are also returned.
// If IncludeMetadata is set to true, any associated metadata of the vectors is also returned, if any.
func (ix *Index) Fetch(f Fetch) (vectors []Vector, err error) {
	return ix.fetchInternal(f, defaultNamespace)
}

func (ix *Index) fetchInternal(f Fetch, ns string) (vectors []Vector, err error) {
	data, err := ix.sendJson(buildPath(fetchPath, ns), f)
	if err != nil {
		return
	}
	vectors, err = parseResponse[[]Vector](data)
	return
}
