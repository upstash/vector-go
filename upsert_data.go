package vector

const upsertDataPath = "/upsert-data"

// UpsertData updates or inserts a vector to the default namespace of the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided while upserting the vector.
func (ix *Index) UpsertData(u UpsertData) (err error) {
	return ix.upsertDataInternal(u, "")
}

// UpsertDataMany updates or inserts some vectors to the default namespace of the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided for each vector.
func (ix *Index) UpsertDataMany(u []UpsertData) (err error) {
	return ix.upsertDataManyInternal(u, "")
}

func (ix *Index) upsertDataInternal(u UpsertData, ns string) (err error) {
	data, err := ix.sendJson(buildPath(upsertDataPath, ns), u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

func (ix *Index) upsertDataManyInternal(u []UpsertData, ns string) (err error) {
	data, err := ix.sendJson(buildPath(upsertDataPath, ns), u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
