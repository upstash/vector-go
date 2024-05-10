package vector

const upsertPath = "/upsert"

// Upsert updates or inserts a vector to the default namespace of the index.
// Additional metadata can also be provided while upserting the vector.
func (ix *Index) Upsert(u Upsert) (err error) {
	return ix.upsertInternal(u, defaultNamespace)
}

// UpsertMany updates or inserts some vectors to the default namespace of the index.
// Additional metadata can also be provided for each vector.
func (ix *Index) UpsertMany(u []Upsert) (err error) {
	return ix.upsertManyInternal(u, defaultNamespace)
}

func (ix *Index) upsertInternal(u Upsert, ns string) (err error) {
	data, err := ix.sendJson(buildPath(upsertPath, ns), u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

func (ix *Index) upsertManyInternal(u []Upsert, ns string) (err error) {
	data, err := ix.sendJson(buildPath(upsertPath, ns), u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
