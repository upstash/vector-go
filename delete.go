package vector

const deletePath = "/delete"

func (c *Client) Delete(id string) (ok bool, err error) {
	data, err := c.sendString(deletePath, id)
	if err != nil {
		return
	}

	deleted, err := parseResponse[Deleted](data)
	ok = deleted.Deleted != 0
	return
}

func (c *Client) DeleteMany(ids []string) (count int, err error) {
	data, err := c.sendJson(deletePath, ids)
	if err != nil {
		return
	}

	deleted, err := parseResponse[Deleted](data)
	count = deleted.Deleted
	return
}
