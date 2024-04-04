package vector

const upsertDataPath = "/upsert-data"

// UpsertData updates or inserts a vector to the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided while upserting the vector.
func (ix *Index) UpsertData(u UpsertData) (err error) {
	data, err := ix.sendJson(upsertDataPath, u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

// UpsertDataMany updates or inserts some vectors to the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided for each vector.
func (ix *Index) UpsertDataMany(u []UpsertData) (err error) {
	data, err := ix.sendJson(upsertDataPath, u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
