package vector

const deletePath = "/delete"

func (c *Client) Delete(id string) (ok bool, err error) {
	// workaround for server that does not accept string bodies
	data, err := c.sendJson(deletePath, []string{id})
	if err != nil {
		return
	}

	deleted, err := parseResponse[deleted](data)
	ok = deleted.Deleted != 0
	return
}

func (c *Client) DeleteMany(ids []string) (count int, err error) {
	data, err := c.sendJson(deletePath, ids)
	if err != nil {
		return
	}

	deleted, err := parseResponse[deleted](data)
	count = deleted.Deleted
	return
}
