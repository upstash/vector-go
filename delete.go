package vector

const deletePath = "/delete"

// Delete deletes the vector with the given id in the default namespace and reports whether the vector is deleted.
// If a vector with the given id is not found, Delete returns false.
func (ix *Index) Delete(id string) (ok bool, err error) {
	return ix.deleteInternal(id, "")
}

// DeleteMany deletes the vectors with the given ids in the default namespace and reports how many of them are deleted.
func (ix *Index) DeleteMany(ids []string) (count int, err error) {
	return ix.deleteManyInternal(ids, "")
}

func (ix *Index) deleteInternal(id string, ns string) (ok bool, err error) {
	data, err := ix.sendBytes(buildPath(deletePath, ns), []byte(id))
	if err != nil {
		return
	}

	deleted, err := parseResponse[deleted](data)
	ok = deleted.Deleted != 0
	return
}

func (ix *Index) deleteManyInternal(ids []string, ns string) (count int, err error) {
	data, err := ix.sendJson(buildPath(deletePath, ns), ids)
	if err != nil {
		return
	}

	deleted, err := parseResponse[deleted](data)
	count = deleted.Deleted
	return
}
