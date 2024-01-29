package vector

const queryPath = "/query"

func (c *Client) Query(q Query) (scores []VectorScore, err error) {
	data, err := c.sendJson(queryPath, q)
	if err != nil {
		return
	}
	scores, err = parseResponse[[]VectorScore](data)
	return
}
