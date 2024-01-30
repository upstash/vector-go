# Upstash Vector Go Client

This is the Go client for [Upstash](https://upstash.com/) Vector.

## Documentation

- [**Reference Documentation**](https://upstash.com/docs/vector/overall/getstarted)

## Installation

```bash
go get github.com/upstash/vector-go
```

## Usage

### Initializing the client

There are two pieces of configuration required to use the Upstash vector client: an REST token and REST URL. 
Find your configuration values in the console dashboard at [https://console.upstash.com/](https://console.upstash.com/).

```go
import (
	"github.com/upstash/vector-go"
)

func main() {
	client := vector.NewClient("<UPSTASH_VECTOR_REST_URL>", "<UPSTASH_VECTOR_REST_TOKEN>")
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
	client := vector.NewClientWith(opts)
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

err := client.UpsertMany(upserts)
```

#### Upsert One

```go
err := client.Upsert(vector.Upsert{
    Id:     "2",
    Vector: []float32{1.0, 0.0},
})
```

### Querying Vectors

The query vector must be present and it must have the same dimensions with the
all the other vectors in the index.

When `TopK` is specified, at most that many vectors will be returned.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
scores, err := client.Query(vector.Query{
    Vector:          []float32{0.0, 1.0},
    TopK:            2,
    IncludeVectors:  false,
    IncludeMetadata: false,
})
```

### Fetching Vectors

Vectors can be fetched individually by providing the unique vector ids.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
vectors, err := client.Fetch(vector.Fetch{
    Ids: []string{"0", "1"},
    IncludeVectors: false,
    IncludeMetadata: false,
})
```

### Deleting Vectors

Vectors can be deleted from the index.

#### Delete many

```go
count, err := client.DeleteMany([]string{"0", "999"})
```

#### Delete One

```go
ok, err := client.Delete("2")
```

### Scanning the Vectors

All or some of the vectors in the index can scanned by fetching range of vectors.

While starting the scan, the initial cursor value of `"0"` should be used.

When `IncludeVectors` is `true`, the response will contain the vector values.

When `IncludeMetadata` is `true`, the response will contain the metadata of the
vectors, if any.

```go
vectors, err := client.Range(vector.Range{
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
    vectors, err = client.Range(vector.Range{
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
err := client.Reset()
```

### Getting Index Information

```go
info, err := client.Info()
```