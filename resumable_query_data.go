package vector

const resumableQueryDataPath = "/resumable-query-data"

// ResumableQueryData starts a resumable query and returns the first page of the
// result of the query for the given text data in the default namespace.
// Then, next pages of the query results can be fetched over the returned handle.
// After all the needed pages of the results are fetched, it is recommended
// to close to handle to release the acquired resources.
func (ix *Index) ResumableQueryData(q ResumableQueryData) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	return ix.resumableQueryDataInternal(q, defaultNamespace)
}

func (ix *Index) resumableQueryDataInternal(q ResumableQueryData, ns string) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	data, err := ix.sendJson(buildPath(resumableQueryDataPath, ns), q)
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
