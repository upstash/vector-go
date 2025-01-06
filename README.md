# Upstash Vector Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/upstash/vector-go.svg)](https://pkg.go.dev/github.com/upstash/vector-go)

> [!NOTE]  
> **This project is in GA Stage.**
>
> The Upstash Professional Support fully covers this project. It receives regular updates, and bug fixes. 
> The Upstash team is committed to maintaining and improving its functionality.

[Upstash](https://upstash.com/) Vector is a serverless vector database designed for working with vector embeddings.

This is the HTTP-based Go client for [Upstash](https://upstash.com/) Vector.

## Documentation

- [**Reference Documentation**](https://upstash.com/docs/vector/overall/getstarted)

## Installation

Use `go get` to install the Upstash Vector package:
```bash
go get github.com/upstash/vector-go
```

Import the Upstash Vector package in your project:

```go
import "github.com/upstash/vector-go"
```

## Usage

In order to use this client, head out to [Upstash Console](https://console.upstash.com) and create a vector database.

### Initializing the client

The REST token and REST URL configurations are required to initialize an Upstash Vector index client.
Find your configuration values in the console dashboard at [Upstash Console](https://console.upstash.com/).

```go
import (
	"github.com/upstash/vector-go"
)

func main() {
	index := vector.NewIndex("<UPSTASH_VECTOR_REST_URL>", "<UPSTASH_VECTOR_REST_TOKEN>")
}
```

Alternatively, you can set following environment variables:

```shell
UPSTASH_VECTOR_REST_URL="your_rest_url"
UPSTASH_VECTOR_REST_TOKEN="your_rest_token"
```

and then create index client by using:

```go
import (
	"github.com/upstash/vector-go"
)

func main() {
	index := vector.NewIndexFromEnv()
}
```

#### Using a custom HTTP client

By default, `http.DefaultClient` will be used for doing requests. It is possible
to use custom HTTP client, by passing it in the options while constructing
the client.

```go
import (
	"net/http"

	"github.com/upstash/vector-go"
)

func main() {
	opts := vector.Options{
		Url:    "<UPSTASH_VECTOR_REST_URL>",
		Token:  "<UPSTASH_VECTOR_REST_TOKEN>",
		Client: &http.Client{},
	}
	index := vector.NewIndexWith(opts)
}
```

## Index operations

Upstash vector indexes support operations for working with vector data using operations such as upsert, query, fetch, and delete.

```go
import (
	"github.com/upstash/vector-go"
)

func main() {
	index := vector.NewIndex("<UPSTASH_VECTOR_REST_URL>", "<UPSTASH_VECTOR_REST_TOKEN>")
}
```

Upstash Vector allows you to partition a single index into multiple isolated namespaces.

You can specify a namespace for an index client with `Namespace(ns string)` function.
When you create a `Namespace` client, all index operations executed through this client become associated with the specified namespace.

By default, the `Index` client is associated with the default namespace.

```go
import (
	"github.com/upstash/vector-go"
)

func main() {
	index := vector.NewIndex("<UPSTASH_VECTOR_REST_URL>", "<UPSTASH_VECTOR_REST_TOKEN>")
	
	// Returns a new Namespace client associated with the given namespace
	ns := index.Namespace("<NAMESPACE>")
```

### Upserting Vectors

Upsert can be used to insert new vectors into index or to update
existing vectors.

#### Upsert Many

**Dense Indexes**

```go
err := index.UpsertMany([]vector.Upsert{
	{
		Id:     "0",
		Vector: []float32{0.6, 0.8},
	},
	{
		Id:       "1",
		Vector:   []float32{0.0, 1.0},
		Metadata: map[string]any{"foo": "bar"}, // optional metadata
		Data:     "vector data",                // optional data
	},
})
```

**Sparse Indexes**
```go
err := index.UpsertMany([]vector.Upsert{
	{
		Id: "0",
		SparseVector: &vector.SparseVector{
			Indices: []int32{0, 1},
			Values:  []float32{0.5, 0.6},
		},
	},
	{
		Id: "1",
		SparseVector: &vector.SparseVector{
			Indices: []int32{5},
			Values:  []float32{0.1},
		},
		Metadata: map[string]any{"foo": "bar"}, // optional metadata
		Data:     "vector data",                // optional data
	},
})
```

**Hybrid Indexes**
```go
err = index.UpsertMany([]vector.Upsert{
	{
		Id:     "0",
		Vector: []float32{0.6, 0.8},
		SparseVector: &vector.SparseVector{
			Indices: []int32{0, 1},
			Values:  []float32{0.5, 0.6},
		},
	},
	{
		Id:     "1",
		Vector: []float32{0.0, 1.0},
		SparseVector: &vector.SparseVector{
			Indices: []int32{5},
			Values:  []float32{0.1},
		},
		Metadata: map[string]any{"foo": "bar"}, // optional metadata
		Data:     "vector data",                // optional data
	},
})
```

#### Upsert One

**Dense Indexes**
```go
err := index.Upsert(vector.Upsert{
	Id:     "2",
	Vector: []float32{1.0, 0.0},
})
```

**Sparse Indexes**
```go
err = index.Upsert(vector.Upsert{
	Id: "2",
	SparseVector: &vector.SparseVector{
		Indices: []int32{0, 1},
		Values:  []float32{0.5, 0.6},
	},
})
```

**Hybrid Indexes**
```go
err = index.Upsert(vector.Upsert{
	Id:     "2",
	Vector: []float32{1.0, 0.0},
	SparseVector: &vector.SparseVector{
		Indices: []int32{0, 1},
		Values:  []float32{0.5, 0.6},
	},
})
```

### Upserting with Raw Data

If the vector index is created with an embedding model, it can be populated using the raw data
without explicitly converting it to an embedding. Upstash server will create the embedding 
and index the generated vectors.

Upsert can be used to insert new vectors into index or to update
existing vectors.

#### Upsert Many

```go
err := index.UpsertDataMany([]vector.UpsertData{
	{
		Id:   "0",
		Data: "Capital of Turkey is Ankara.",
	},
	{
		Id:       "1",
		Data:     "Capital of Japan is Tokyo.",
		Metadata: map[string]any{"foo": "bar"}, // optional metadata
	},
})
```

#### Upsert One

```go
err := index.UpsertData(vector.UpsertData{
	Id:   "2",
	Data: "Capital of Turkey is Ankara.",
})
```

### Querying Vectors

When `TopK` is specified, at most that many vectors will be returned.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

When `IncludeData` is `true`, the response will contain the data of the vectors, if any.

**Dense Indexes**
```go
scores, err := index.Query(vector.Query{
	Vector:          []float32{0.0, 1.0},
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
```

**Sparse Indexes**
```go
scores, err := index.Query(vector.Query{
	SparseVector: &vector.SparseVector{
		Indices: []int32{0, 1},
		Values:  []float32{0.5, 0.5},
	},
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
```

**Hybrid Indexes**
```go
scores, err := index.Query(vector.Query{
	Vector: []float32{0.0, 1.0},
	SparseVector: &vector.SparseVector{
		Indices: []int32{0, 1},
		Values:  []float32{0.5, 0.5},
	},
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
```

Additionally, a metadata filter can be specified in queries. When `Filter` is given, the response will contain 
only the values whose metadata matches the given filter. See [Metadata Filtering](https://upstash.com/docs/vector/features/metadatafiltering) 
docs for more information.

```go
scores, err := index.Query(vector.Query{
	..., 
	Filter: `foo = 'bar'`
})
```

### Querying with Raw Data

If the vector index is created with an embedding model, a query can be executed using the raw data
without explicitly converting it to an embedding. Upstash server will create the embedding and run the query.

When `TopK` is specified, at most that many vectors will be returned.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

When `IncludeData` is `true`, the response will contain the data of the vectors, if any.

```go
scores, err := index.QueryData(vector.QueryData{
	Data:            "Where is the capital of Turkey?",
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
	Filter:          `foo = 'bar'`,
})
```

### Resumable Querying Vectors

With a similar interface to query and query data, query results
can be fetched page by page with resumable queries.

#### Resumalbe Query

When a resumable query is started, it returns the first
page of the query results, and returns a handle that
can be used to fetch next pages. When enough query results
are fetched, handle can be closed to release the resources
acquired in the index to facilitate the resumable query.

**Dense Indexes**
```go
scores, handle, err := index.ResumableQuery(vector.ResumableQuery{
	Vector:          []float32{0.0, 1.0},
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
defer handle.Close()

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 3,
})

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 5,
})
```

**Sparse Indexes**
```go
scores, handle, err := index.ResumableQuery(vector.ResumableQuery{
	SparseVector: &vector.SparseVector{
		Indices: []int32{0, 1},
		Values:  []float32{0.5, 0.5},
	},
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
defer handle.Close()

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 3,
})

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 5,
})
```

**Hybrid Indexes**
```go
scores, handle, err := index.ResumableQuery(vector.ResumableQuery{
	Vector: []float32{0.0, 1.0},
	SparseVector: &vector.SparseVector{
		Indices: []int32{0, 1},
		Values:  []float32{0.5, 0.5},
	},
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
defer handle.Close()

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 3,
})

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 5,
})
```

#### Resumable Query with Data

If the vector index is created with an embedding model, a resumable query can be started using the raw data
without explicitly converting it to an embedding. Upstash server will create the embedding and start the query.

```go
scores, handle, err := index.ResumableQueryData(vector.ResumableQueryData{
	Data:            "Where is the capital of Turkey?",
	TopK:            2,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
defer handle.Close()

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 3,
})

scores, err = handle.Next(vector.ResumableQueryNext{
	AdditionalK: 5,
})
```

### Fetching Vectors

Vectors can be fetched individually by providing the unique vector ids.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

When `IncludeData` is `true`, the response will contain the data of the vectors, if any.

```go
vectors, err := index.Fetch(vector.Fetch{
	Ids:             []string{"0", "1"},
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})
```

### Deleting Vectors

Vectors can be deleted from the index.

#### Delete many

```go
count, err := index.DeleteMany([]string{"0", "999"})
```

#### Delete One

```go
ok, err := index.Delete("2")
```

### Scanning Vectors

All or some of the vectors in the index can scanned by fetching range of vectors.

While starting the scan, the initial cursor value of `"0"` should be used.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

When `IncludeData` is `true`, the response will contain the data of the vectors, if any.

```go
vectors, err := index.Range(vector.Range{
	Cursor:          "0",
	Limit:           10,
	IncludeVectors:  false,
	IncludeMetadata: false,
	IncludeData:     false,
})

for vectors.NextCursor != "" {
	for _, v := range vectors.Vectors {
		// process individual vectors
	}

	// Fetch the next range batch
	vectors, err = index.Range(vector.Range{
		Cursor:          vectors.NextCursor,
		Limit:           10,
		IncludeVectors:  false,
		IncludeMetadata: false,
		IncludeData:     false,
	})
}
```

### Updating Vectors

Any combination of vector value, sparse vector value, data, or metadata can be updated.

```go
ok, err := index.Update(vector.Update{
	Id:       "id",
	Metadata: map[string]any{"new": "metadata"},
})
```

### Resetting the Index

Reset will delete all the vectors and reset the index to its initial state.

```go
err := index.Reset()
```

### Getting Index Information

```go
info, err := index.Info()
```

### List Namespaces

All the names of active namespaces can be listed.

```go
namespaces, err := index.ListNamespaces()
for _, ns : range namespaces {
	fmt.Println(ns)
}
```

### Delete Namespaces

A namespace can be deleted entirely if it exists.
The default namespaces cannot be deleted.

```go
err := index.Namespace("ns").DeleteNamespace()
```
