package fsutil

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnionFS(t *testing.T) {
	t.Parallel()

	fs1 := fstest.MapFS{
		"f1": &fstest.MapFile{Data: []byte("fs1")},
		"f2": &fstest.MapFile{Data: []byte("fs1")},
	}
	fs2 := fstest.MapFS{
		"f1": &fstest.MapFile{Data: []byte("fs2")},
		"f3": &fstest.MapFile{Data: []byte("fs2")},
	}

	u := UnionFS{fs1, fs2}

	fstest.TestFS(u, "f1", "f2", "f3")
	assertContent(t, u, "f1", "fs1")
	assertContent(t, u, "f2", "fs1")
	assertContent(t, u, "f3", "fs2")
	assertNotExists(t, u, "f4")

	t.Run("ReadDir", func(t *testing.T) {
		files, err := fs.ReadDir(u, ".")
		require.NoError(t, err)
		var names []string
		for _, f := range files {
			names = append(names, f.Name())
		}
		want := []string{"f1", "f2", "f3"}
		assert.Equal(t, want, names)
	})
}

func TestUnionFSSub(t *testing.T) {
	t.Parallel()

	fs1 := fstest.MapFS{
		"a/f1": &fstest.MapFile{Data: []byte("fs1")},
		"a/f2": &fstest.MapFile{Data: []byte("fs1")},
		"b/f3": &fstest.MapFile{Data: []byte("fs1")},
	}
	fs2 := fstest.MapFS{
		"a/f1": &fstest.MapFile{Data: []byte("fs2")},
		"a/f4": &fstest.MapFile{Data: []byte("fs2")},
		"b/f5": &fstest.MapFile{Data: []byte("fs1")},
	}

	u := UnionFS{fs1, fs2}

	a, err := fs.Sub(u, "a")
	require.NoError(t, err)
	fstest.TestFS(a, "f1", "f2", "f4")
	assertNotExists(t, a, "f3")
	assertNotExists(t, a, "f5")

	b, err := fs.Sub(u, "b")
	require.NoError(t, err)
	fstest.TestFS(b, "f3", "f5")
	assertNotExists(t, b, "f1")
	assertNotExists(t, b, "f2")
	assertNotExists(t, b, "f4")

	c, err := fs.Sub(u, "c")
	require.NoError(t, err)
	fstest.TestFS(c)
}
