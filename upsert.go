package vector

const upsertPath = "/upsert"

// Upsert updates or inserts a vector to the index.
// Additional metadata can also be provided while upserting the vector.
func (ix *Index) Upsert(u Upsert) (err error) {
	data, err := ix.sendJson(upsertPath, u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

// UpsertMany updates or inserts some vectors to the index.
// Additional metadata can also be provided for each vector.
func (ix *Index) UpsertMany(u []Upsert) (err error) {
	data, err := ix.sendJson(upsertPath, u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
