package vector

const infoPath = "/info"

// Info returns some information about the index; such as
// vector count, vectors pending for indexing, size of the
// index on disk in bytes, dimension of the index, and
// the name of the similarity function used.
func (c *Client) Info() (info IndexInfo, err error) {
	data, err := c.sendJson(infoPath, nil)
	if err != nil {
		return
	}
	info, err = parseResponse[IndexInfo](data)
	return
}
