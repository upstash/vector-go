package vector

const upsertPath = "/upsert"

// Upsert updates or inserts a vector to a namespace of the index.
// Additional metadata can also be provided while upserting the vector.
// If namespace is not specified, the default namespace is used.
func (ix *Index) Upsert(u Upsert) (err error) {
	data, err := ix.sendJson(upsertPath, u, true)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

// UpsertMany updates or inserts some vectors to a namespace of the index.
// Additional metadata can also be provided for each vector.
// If namespace is not specified, the default namespace is used.
func (ix *Index) UpsertMany(u []Upsert) (err error) {
	data, err := ix.sendJson(upsertPath, u, true)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
