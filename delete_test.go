package vector

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomString() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

func TestDelete(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace_"+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeDense, ns)
			require.NoError(t, err)
			id := randomString()
			err = client.Upsert(Upsert{
				Id:     id,
				Vector: []float32{0, 1},
			})
			require.NoError(t, err)

			t.Run("existing id", func(t *testing.T) {
				ok, err := client.Delete(id)
				require.NoError(t, err)
				require.True(t, ok)
			})

			t.Run("non existing id", func(t *testing.T) {
				ok, err := client.Delete(randomString())
				require.NoError(t, err)
				require.False(t, ok)
			})
		})
	}
}

func TestDeleteMany(t *testing.T) {
	for _, ns := range namespaces {
		t.Run("namespace "+ns, func(t *testing.T) {
			client, err := newTestClient(testClientTypeDense, ns)
			require.NoError(t, err)

			id0 := randomString()
			id1 := randomString()
			id2 := randomString()
			err = client.UpsertMany([]Upsert{
				{
					Id:     id0,
					Vector: []float32{0, 1},
				},
				{
					Id:     id1,
					Vector: []float32{5, 10},
				},
			})
			require.NoError(t, err)

			count, err := client.DeleteMany([]string{id0, id1, id2})
			require.NoError(t, err)
			require.Equal(t, 2, count)
		})
	}
}
