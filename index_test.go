package vector

import (
	"errors"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	namespaces      = [...]string{defaultNamespace, "ns"}
	testClientTypes = [...]testClientType{testClientTypeDense, testClientTypeSparse, testClientTypeHybrid}
)

func init() {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(err)
	}
}

// We are using the same internal methods for indexes and namespaces,
// the only difference being index using the default namespace name.
// However, we are left with using only the namespace or index in our
// tests if we were to use these types directly(or duplicate all the tests).
// So there is no way of verifying that the namespace or index
// uses the correct internal methods, apart from the manual inspection.
// Instead, this class will call the appropriate methods over the Index
// for the default namespace, and over the Namespace for others.
type testClient struct {
	index         *Index
	namespace     *Namespace
	namespaceName string
}

func (tc *testClient) Upsert(u Upsert) (err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Upsert(u)
	} else {
		return tc.namespace.Upsert(u)
	}
}

func (tc *testClient) UpsertMany(u []Upsert) (err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.UpsertMany(u)
	} else {
		return tc.namespace.UpsertMany(u)
	}
}

func (tc *testClient) UpsertData(u UpsertData) (err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.UpsertData(u)
	} else {
		return tc.namespace.UpsertData(u)
	}
}

func (tc *testClient) UpsertDataMany(u []UpsertData) (err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.UpsertDataMany(u)
	} else {
		return tc.namespace.UpsertDataMany(u)
	}
}

func (tc *testClient) Fetch(f Fetch) (vectors []Vector, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Fetch(f)
	} else {
		return tc.namespace.Fetch(f)
	}
}

func (tc *testClient) QueryData(q QueryData) (scores []VectorScore, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.QueryData(q)
	} else {
		return tc.namespace.QueryData(q)
	}
}

func (tc *testClient) Query(q Query) (scores []VectorScore, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Query(q)
	} else {
		return tc.namespace.Query(q)
	}
}

func (tc *testClient) Range(r Range) (vectors RangeVectors, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Range(r)
	} else {
		return tc.namespace.Range(r)
	}
}

func (tc *testClient) Delete(id string) (ok bool, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Delete(id)
	} else {
		return tc.namespace.Delete(id)
	}
}

func (tc *testClient) DeleteMany(ids []string) (count int, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.DeleteMany(ids)
	} else {
		return tc.namespace.DeleteMany(ids)
	}
}

func (tc *testClient) Reset() (err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Reset()
	} else {
		return tc.namespace.Reset()
	}
}

func (tc *testClient) Update(u Update) (ok bool, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.Update(u)
	} else {
		return tc.namespace.Update(u)
	}
}

func (tc *testClient) ResumableQuery(q ResumableQuery) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.ResumableQuery(q)
	} else {
		return tc.namespace.ResumableQuery(q)
	}
}

func (tc *testClient) ResumableQueryData(q ResumableQueryData) (scores []VectorScore, handle *ResumableQueryHandle, err error) {
	if tc.namespaceName == defaultNamespace {
		return tc.index.ResumableQueryData(q)
	} else {
		return tc.namespace.ResumableQueryData(q)
	}
}

func (tc *testClient) Info() (info IndexInfo, err error) {
	return tc.index.Info()
}

type testClientType string

const (
	testClientTypeDense           testClientType = "dense"
	testClientTypeSparse          testClientType = "sparse"
	testClientTypeHybrid          testClientType = "hybrid"
	testClientTypeDenseEmbedding  testClientType = "dense_embedding"
	testClientTypeHybridEmbedding testClientType = "hybrid_embedding"
)

func newTestClient(t testClientType, ns string) (*testClient, error) {
	var urlEnv, tokenEnv string
	switch t {
	case testClientTypeDense:
		urlEnv = UrlEnvProperty
		tokenEnv = TokenEnvProperty
	case testClientTypeSparse:
		urlEnv = "SPARSE_" + UrlEnvProperty
		tokenEnv = "SPARSE_" + TokenEnvProperty
	case testClientTypeHybrid:
		urlEnv = "HYBRID_" + UrlEnvProperty
		tokenEnv = "HYBRID_" + TokenEnvProperty
	case testClientTypeDenseEmbedding:
		urlEnv = "EMBEDDING_" + UrlEnvProperty
		tokenEnv = "EMBEDDING_" + TokenEnvProperty
	case testClientTypeHybridEmbedding:
		urlEnv = "HYBRID_EMBEDDING_" + UrlEnvProperty
		tokenEnv = "HYBRID_EMBEDDING_" + TokenEnvProperty
	}

	opts := Options{
		Url:   os.Getenv(urlEnv),
		Token: os.Getenv(tokenEnv),
	}

	index := NewIndexWith(opts)

	for _, n := range namespaces {
		err := index.Namespace(n).Reset()
		if err != nil {
			return nil, err
		}
	}

	namespace := index.Namespace(ns)

	return &testClient{
		index:         index,
		namespace:     namespace,
		namespaceName: ns,
	}, nil
}

func randomVectors(tcType testClientType) ([]float32, *SparseVector) {
	switch tcType {
	case testClientTypeDense:
		return []float32{rand.Float32(), rand.Float32()}, nil
	case testClientTypeSparse:
		return nil, &SparseVector{
			Indices: []int32{rand.Int31(), rand.Int31()},
			Values:  []float32{rand.Float32(), rand.Float32()},
		}
	default:
		return []float32{rand.Float32(), rand.Float32()}, &SparseVector{
			Indices: []int32{rand.Int31(), rand.Int31()},
			Values:  []float32{rand.Float32(), rand.Float32()},
		}
	}
}

func TestNewTestClient(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeDense, ns)
			require.NoError(t, err)

			_, err = client.Info()
			require.NoError(t, err)
		})
	}
}

func TestNewIndex(t *testing.T) {
	c := NewIndex(os.Getenv(UrlEnvProperty), os.Getenv(TokenEnvProperty))

	_, err := c.Info()
	require.NoError(t, err)
}

func TestNewIndexWith(t *testing.T) {
	opts := Options{
		Url:    os.Getenv(UrlEnvProperty),
		Token:  os.Getenv(TokenEnvProperty),
		Client: &http.Client{},
	}

	c := NewIndexWith(opts)

	_, err := c.Info()
	require.NoError(t, err)
}

func TestNewIndexFromEnv(t *testing.T) {
	c := NewIndexFromEnv()

	_, err := c.Info()
	require.NoError(t, err)
}
