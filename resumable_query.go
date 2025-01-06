package vector

const (
	resumableQueryPath    = "/resumable-query"
	resumableQueryNexPath = "/resumable-query-next"
	resumableQueryEndPath = "/resumable-query-end"
)

type ResumableQueryHandle struct {
	index *Index
	uuid  string
}

// Next fetches the next page of the query result.
func (h *ResumableQueryHandle) Next(n ResumableQueryNext) (scores []VectorScore, err error) {
	nn := resumableQueryNext{
		ResumableQueryNext: n,
		UUID:               h.uuid,
	}

	data, err := h.index.sendJson(buildPath(resumableQueryNexPath, defaultNamespace), nn)
	if err != nil {
		return
	}

	scores, err = parseResponse[[]VectorScore](data)
	return
}

// Close stops the resumable query and releases the acquired resources.
func (h *ResumableQueryHandle) Close() (err error) {
	e := resumableQueryEnd{UUID: h.uuid}
	data, err := h.index.sendJson(buildPath(resumableQueryEndPath, defaultNamespace), e)
	if err != nil {
		return
	}

	_, err = parseResponse[string](data)
	return
}

// ResumableQuery starts a resumable query and returns the first page of the
// result of the query for the given vector in the default namespace.
// Then, next pages of the query results can be fetched over the returned handle.
// After all the needed pages of the results are fetched, it is recommended
// to close to handle to release the acquired resources.
func (ix *Index) ResumableQuery(q ResumableQuery) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	return ix.resumableQueryInternal(q, defaultNamespace)
}

func (ix *Index) resumableQueryInternal(q ResumableQuery, ns string) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	data, err := ix.sendJson(buildPath(resumableQueryPath, ns), q)
	if err != nil {
		return
	}

	start, err := parseResponse[resumableQueryStart](data)
	if err != nil {
		return
	}

	scores = start.Scores
	handle = &ResumableQueryHandle{
		index: ix,
		uuid:  start.UUID,
	}
	return
}
