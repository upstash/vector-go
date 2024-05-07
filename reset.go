package vector

const resetPath = "/reset"

// Reset deletes all the vectors in a namespace of the index and resets it to initial state.
// When not specified, the default namespace is used.
// Use Namespace(ns string) method to specify a namespace for the client.
func (ix *Index) Reset() (err error) {
	data, err := ix.send(resetPath, nil, true)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
