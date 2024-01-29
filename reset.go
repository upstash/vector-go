package vector

const resetPath = "/reset"

func (c *Client) Reset() (err error) {
	data, err := c.send(resetPath, nil)
	if err != nil {
		return
	}
	_, err = parseResponse[string](data)
	return
}
