package vector

const resetPath = "/reset"

// Reset deletes all the vectors in a namespace of the index and resets it to initial state.
// If namespace is not specified, the default namespace is used.
func (ix *Index) Reset() (err error) {
	data, err := ix.send(resetPath, nil, true)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
