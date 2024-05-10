package vector

const rangePath = "/range"

// Range returns a range of vectors, starting with r.Cursor (inclusive),
// until the end of the vectors in the index or until the given q.Limit.
// The initial cursor should be set to "0", and subsequent calls to
// Range might use the next cursor returned in the response.
// When r.IncludeVectors is true, values of the vectors are also returned.
// When r.IncludeMetadata is true, metadata of the vectors are also returned, if any.
func (ix *Index) Range(r Range) (vectors RangeVectors, err error) {
	return ix.rangeInternal(r, "")
}

func (ix *Index) rangeInternal(r Range, ns string) (vectors RangeVectors, err error) {
	data, err := ix.sendJson(buildPath(rangePath, ns), r)
	if err != nil {
		return
	}
	vectors, err = parseResponse[RangeVectors](data)
	return
}
