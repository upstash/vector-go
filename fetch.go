package vector

const fetchPath = "/fetch"

func (c *Client) Fetch(f Fetch) (vectors []Vector, err error) {
	data, err := c.sendJson(fetchPath, f)
	if err != nil {
		return
	}
	vectors, err = parseResponse[[]Vector](data)
	return
}
