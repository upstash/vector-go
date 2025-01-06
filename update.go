package vector

const updatePath = "/update"

// Update updates a vector value, data, or metadata for the given id
// for the default namespace of the index and reports whether the vector is updated.
// If a vector with the given id is not found, Update returns false.
func (ix *Index) Update(u Update) (ok bool, err error) {
	return ix.updateInternal(u, defaultNamespace)
}

func (ix *Index) updateInternal(u Update, ns string) (ok bool, err error) {
	data, err := ix.sendJson(buildPath(updatePath, ns), u)
	if err != nil {
		return
	}

	res, err := parseResponse[updated](data)
	if err != nil {
		return
	}

	ok = res.Updated == 1
	return
}
