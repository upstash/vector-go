package vector

const resetPath = "/reset"

// Reset deletes all the vectors in the default namespace of the index and resets it to initial state.
func (ix *Index) Reset() (err error) {
	return ix.resetInternal("")
}

func (ix *Index) resetInternal(ns string) (err error) {
	data, err := ix.send(buildPath(resetPath, ns), nil)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
