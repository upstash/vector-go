package vector

const rangePath = "/range"

func (c *Client) Range(r Range) (vectors RangeVectors, err error) {
	data, err := c.sendJson(rangePath, r)
	if err != nil {
		return
	}
	vectors, err = parseResponse[RangeVectors](data)
	return
}
