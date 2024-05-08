package vector

const deletePath = "/delete"

// Delete deletes the vector in a namespace with the given id and reports whether the vector is deleted.
// If a vector with the given id is not found, Delete returns false.
// If namespace is not specified, the default namespace is used.
func (ix *Index) Delete(id string) (ok bool, err error) {
	data, err := ix.sendBytes(deletePath, []byte(id), true)
	if err != nil {
		return
	}

	deleted, err := parseResponse[deleted](data)
	ok = deleted.Deleted != 0
	return
}

// DeleteMany deletes the vectors in a namespace with the given ids and reports how many of them are deleted.
// If namespace is not specified, the default namespace is used.
func (ix *Index) DeleteMany(ids []string) (count int, err error) {
	data, err := ix.sendJson(deletePath, ids, true)
	if err != nil {
		return
	}

	deleted, err := parseResponse[deleted](data)
	count = deleted.Deleted
	return
}
