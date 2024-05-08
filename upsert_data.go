package vector

const upsertDataPath = "/upsert-data"

// UpsertData updates or inserts a vector to a namespace of the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided while upserting the vector.
// If namespace is not specified, the default namespace is used.
func (ix *Index) UpsertData(u UpsertData) (err error) {
	data, err := ix.sendJson(upsertDataPath, u, true)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

// UpsertDataMany updates or inserts some vectors to a namespace of the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided for each vector.
// If namespace is not specified, the default namespace is used.
func (ix *Index) UpsertDataMany(u []UpsertData) (err error) {
	data, err := ix.sendJson(upsertDataPath, u, true)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
