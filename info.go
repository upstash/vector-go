package vector

const infoPath = "/info"

// Info returns some information about the index, including:
//   - Total number of vectors across all namespaces
//   - Total number of vectors waiting to be indexed across all namespaces
//   - Total size of the index on disk in bytes
//   - Vector dimension
//   - Similarity function used
//   - per-namespace vector and pending vector counts
func (ix *Index) Info() (info IndexInfo, err error) {
	data, err := ix.sendJson(infoPath, nil, false)
	if err != nil {
		return
	}
	info, err = parseResponse[IndexInfo](data)
	return
}
