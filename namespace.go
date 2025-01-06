package vector

const deleteNamespacePath = "/delete-namespace"
const listNamespacesPath = "/list-namespaces"

type Namespace struct {
	index *Index
	ns    string
}

// Namespace returns a new client associated with the given namespace.
func (ix *Index) Namespace(namespace string) (i *Namespace) {
	return &Namespace{
		index: ix,
		ns:    namespace,
	}
}

// ListNamespaces returns the list of names of namespaces.
func (ix *Index) ListNamespaces() (namespaces []string, err error) {
	data, err := ix.sendJson(listNamespacesPath, nil)
	if err != nil {
		return
	}
	namespaces, err = parseResponse[[]string](data)
	return
}

// DeleteNamespace deletes the given namespace of index if it exists.
func (ns *Namespace) DeleteNamespace() error {
	_, err := ns.index.sendBytes(buildPath(deleteNamespacePath, ns.ns), nil)
	return err
}

// Upsert updates or inserts a vector to the namespace of the index.
// Additional metadata can also be provided while upserting the vector.
func (ns *Namespace) Upsert(u Upsert) (err error) {
	return ns.index.upsertInternal(u, ns.ns)
}

// UpsertMany updates or inserts some vectors to the default namespace of the index.
// Additional metadata can also be provided for each vector.
func (ns *Namespace) UpsertMany(u []Upsert) (err error) {
	return ns.index.upsertManyInternal(u, ns.ns)
}

// UpsertData updates or inserts a vector to the namespace of the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided while upserting the vector.
func (ns *Namespace) UpsertData(u UpsertData) (err error) {
	return ns.index.upsertDataInternal(u, ns.ns)
}

// UpsertDataMany updates or inserts some vectors to the default namespace of the index
// by converting given raw data to an embedding on the server.
// Additional metadata can also be provided for each vector.
func (ns *Namespace) UpsertDataMany(u []UpsertData) (err error) {
	return ns.index.upsertDataManyInternal(u, ns.ns)
}

// Fetch fetches one or more vectors in the namespace with the ids passed into f.
// If IncludeVectors is set to true, the vector values are also returned.
// If IncludeMetadata is set to true, any associated metadata of the vectors is also returned, if any.
func (ns *Namespace) Fetch(f Fetch) (vectors []Vector, err error) {
	return ns.index.fetchInternal(f, ns.ns)
}

// QueryData returns the result of the query for the given data by converting it to an embedding on the server.
// When q.TopK is specified, the result will contain at most q.TopK many vectors.
// The returned list will contain vectors sorted in descending order of score,
// which correlates with the similarity of the vectors to the given query vector.
// When q.IncludeVectors is true, values of the vectors are also returned.
// When q.IncludeMetadata is true, metadata of the vectors are also returned, if any.
func (ns *Namespace) QueryData(q QueryData) (scores []VectorScore, err error) {
	return ns.index.queryDataInternal(q, ns.ns)
}

// Query returns the result of the query for the given vector in the namespace.
// When q.TopK is specified, the result will contain at most q.TopK many vectors.
// The returned list will contain vectors sorted in descending order of score,
// which correlates with the similarity of the vectors to the given query vector.
// When q.IncludeVectors is true, values of the vectors are also returned.
// When q.IncludeMetadata is true, metadata of the vectors are also returned, if any.
func (ns *Namespace) Query(q Query) (scores []VectorScore, err error) {
	return ns.index.queryInternal(q, ns.ns)
}

// Range returns a range of vectors, starting with r.Cursor (inclusive),
// until the end of the vectors in the index or until the given q.Limit.
// The initial cursor should be set to "0", and subsequent calls to
// Range might use the next cursor returned in the response.
// When r.IncludeVectors is true, values of the vectors are also returned.
// When r.IncludeMetadata is true, metadata of the vectors are also returned, if any.
func (ns *Namespace) Range(r Range) (vectors RangeVectors, err error) {
	return ns.index.rangeInternal(r, ns.ns)
}

// Delete deletes the vector with the given id in the namespace and reports whether the vector is deleted.
// If a vector with the given id is not found, Delete returns false.
func (ns *Namespace) Delete(id string) (ok bool, err error) {
	return ns.index.deleteInternal(id, ns.ns)
}

// DeleteMany deletes the vectors with the given ids in the namespace and reports how many of them are deleted.
func (ns *Namespace) DeleteMany(ids []string) (count int, err error) {
	return ns.index.deleteManyInternal(ids, ns.ns)
}

// Reset deletes all the vectors in the namespace of the index and resets it to initial state.
func (ns *Namespace) Reset() (err error) {
	return ns.index.resetInternal(ns.ns)
}

// Update updates a vector value, data, or metadata for the given id
// in the namespace and reports whether the vector is updated.
// If a vector with the given id is not found, Update returns false.
func (ns *Namespace) Update(u Update) (ok bool, err error) {
	return ns.index.updateInternal(u, ns.ns)
}

// ResumableQuery starts a resumable query and returns the first page of the
// result of the query for the given vector in the default namespace.
// Then, next pages of the query results can be fetched over the returned handle.
// After all the needed pages of the results are fetched, it is recommended
// to close to handle to release the acquired resources.
func (ns *Namespace) ResumableQuery(q ResumableQuery) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	return ns.index.resumableQueryInternal(q, ns.ns)
}

// ResumableQueryData starts a resumable query and returns the first page of the
// result of the query for the given text data in the default namespace.
// Then, next pages of the query results can be fetched over the returned handle.
// After all the needed pages of the results are fetched, it is recommended
// to close to handle to release the acquired resources.
func (ns *Namespace) ResumableQueryData(q ResumableQueryData) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	return ns.index.resumableQueryDataInternal(q, ns.ns)
}
