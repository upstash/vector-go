package vector

const resetPath = "/reset"

// Reset deletes all the vectors in the index and resets
// it to initial state.
func (ix *Index) Reset() (err error) {
	data, err := ix.send(resetPath, nil)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
