package vector

const infoPath = "/info"

func (c *Client) Info() (info IndexInfo, err error) {
	data, err := c.sendJson(infoPath, nil)
	if err != nil {
		return
	}
	info, err = parseResponse[IndexInfo](data)
	return
}
