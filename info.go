package vector

const infoPath = "/info"

// Info returns some information about the index; such as
// vector count, vectors pending for indexing, size of the
// index on disk in bytes, dimension of the index, and
// the name of the similarity function used.
func (ix *Index) Info() (info IndexInfo, err error) {
	data, err := ix.sendJson(infoPath, nil, false)
	if err != nil {
		return
	}
	info, err = parseResponse[IndexInfo](data)
	return
}
