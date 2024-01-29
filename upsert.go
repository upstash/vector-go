package vector

const upsertPath = "/upsert"

func (c *Client) Upsert(u Upsert) (err error) {
	data, err := c.sendJson(upsertPath, u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}

func (c *Client) UpsertMany(u []Upsert) (err error) {
	data, err := c.sendJson(upsertPath, u)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
