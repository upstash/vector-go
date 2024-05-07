package vector

import "errors"

const deleteNamespacePath = "/delete-namespace"
const listNamespacesPath = "/list-namespaces"

func (ix *Index) Namespace(ns string) (i *Index) {
	return &Index{
		url:       ix.url,
		token:     ix.token,
		client:    ix.client,
		namespace: ns,
		headers:   ix.headers,
	}
}

// DeleteNamespace deletes the given namespace of index if it exists.
func (ix *Index) DeleteNamespace(ns string) error {
	if ns == "" {
		return errors.New("cannot delete the default namespace")
	}
	pns := ix.namespace
	ix.namespace = ns
	_, err := ix.sendBytes(deleteNamespacePath, nil, true)
	ix.namespace = pns
	return err
}

// ListNamespaces returns the list of names of namespaces for the index.
func (ix *Index) ListNamespaces() (namespaces []string, err error) {
	data, err := ix.sendJson(listNamespacesPath, nil, false)
	if err != nil {
		return
	}
	namespaces, err = parseResponse[[]string](data)
	return
}
