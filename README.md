# Upstash Vector Go Client

> [!NOTE]  
> **This project is in GA Stage.**
>
> The Upstash Professional Support fully covers this project. It receives regular updates, and bug fixes. The Upstash team is committed to maintaining and improving its functionality.

This is the Go client for [Upstash](https://upstash.com/) Vector.

## Documentation

- [**Reference Documentation**](https://upstash.com/docs/vector/overall/getstarted)

## Installation

```bash
go get github.com/upstash/vector-go
```

## Usage

### Initializing the client

There are two pieces of configuration required to use the Upstash Vector index client: an REST token and REST URL. 
Find your configuration values in the console dashboard at [https://console.upstash.com/](https://console.upstash.com/).

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

### Upserting Vectors

All vectors upserted to index must have the same dimensions.

Upsert can be used to insert new vectors into index or to update
existing vectors.

#### Upsert many

```go
upserts := []vector.Upsert{
    {
        Id:     "0",
        Vector: []float32{0.6, 0.8},
    },
    {
        Id:       "1",
        Vector:   []float32{0.0, 1.0},
        Metadata: map[string]any{"foo": "bar"}, // optional metadata
    },
}

err := index.UpsertMany(upserts)
```

#### Upsert One

```go
err := index.Upsert(vector.Upsert{
    Id:     "2",
    Vector: []float32{1.0, 0.0},
})
```

### Upserting with Raw Data

If the vector index is created with a predefined embedding model, it can be populated using the raw data
without explicitly converting it to an embedding. Upstash server will create the embedding 
and index the generated vectors.

Upsert can be used to insert new vectors into index or to update
existing vectors.

#### Upsert many

```go
upserts := []vector.UpsertData{
    {
        Id:     "0",
        Data:   "Capital of Turkey is Ankara.",
    },
    {
        Id:       "1",
        Data:     "Capital of Japan is Tokyo.",
        Metadata: map[string]any{"foo": "bar"}, // optional metadata
    },
}

err := index.UpsertDataMany(upserts)
```

#### Upsert One

```go
err := index.UpsertData(vector.UpsertData{
    Id:     "2",
    Data:   "Capital of Turkey is Ankara.",
})
```

### Querying Vectors

The query vector must be present, and it must have the same dimensions with the
all the other vectors in the index.

When `TopK` is specified, at most that many vectors will be returned.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
scores, err := index.Query(vector.Query{
    Vector:          []float32{0.0, 1.0},
    TopK:            2,
    IncludeVectors:  false,
    IncludeMetadata: false,
})
```

Additionally, a metadata filter can be specified in queries. When `Filter` is given, the response will contain 
only the values whose metadata matches the given filter. See [Metadata Filtering](https://upstash.com/docs/vector/features/metadatafiltering) 
docs for more information.

```go
scores, err := index.Query(vector.Query{
    Vector:          []float32{0.0, 1.0},
    TopK:            2,
    IncludeVectors:  false,
    IncludeMetadata: false,
    Filter: `foo = 'bar'`
})
```

### Querying with Raw Data

If the vector index is created with a predefined embedding model, a query can be executed using the raw data
without explicitly converting it to an embedding. Upstash server will create the embedding and run the query.

When `TopK` is specified, at most that many vectors will be returned.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
scores, err := index.QueryData(vector.QueryData{
    Data:            "Where is the capital of Turkey?",
    TopK:            2,
    IncludeVectors:  false,
    IncludeMetadata: false,
    Filter: `foo = 'bar'`
})
```

### Fetching Vectors

Vectors can be fetched individually by providing the unique vector ids.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
vectors, err := index.Fetch(vector.Fetch{
    Ids: []string{"0", "1"},
    IncludeVectors: false,
    IncludeMetadata: false,
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

### Scanning the Vectors

All or some of the vectors in the index can scanned by fetching range of vectors.

While starting the scan, the initial cursor value of `"0"` should be used.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
vectors, err := index.Range(vector.Range{
    Cursor:          "0",
    Limit:           10,
    IncludeVectors:  false,
    IncludeMetadata: false,
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
    })
}
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
