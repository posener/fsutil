package fsutil

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefixFS(t *testing.T) {
	t.Parallel()

	f := fstest.MapFS{"c/d": &fstest.MapFile{Data: []byte("data")}}
	p := AddPrefix(f, "aa/b")

	assertNotExists(t, p, "c/d")
	assertContent(t, p, "aa/b/c/d", "data")

	fstest.TestFS(p, "aa/b/c/d")

	t.Run("Open", func(t *testing.T) {
		tests := []struct {
			name  string
			isDir bool
		}{
			{name: "aa", isDir: true},
			{name: "aa/b", isDir: true},
			{name: "aa/b/c", isDir: true},
			{name: "aa/b/c/d", isDir: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				f, err := p.Open(tt.name)
				require.NoError(t, err)
				s, err := f.Stat()
				require.NoError(t, err)
				assert.Equal(t, tt.isDir, s.IsDir())
			})
		}
	})

	t.Run("OpenNotExist", func(t *testing.T) {
		tests := []struct {
			name string
		}{
			{name: "a"},
			{name: "b"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := p.Open(tt.name)
				assert.Equal(t, err, fs.ErrNotExist)
			})
		}
	})

	t.Run("Sub", func(t *testing.T) {
		tests := []struct {
			dir  string
			want string
		}{
			{dir: "aa", want: "b/c/d"},
			{dir: "aa/b", want: "c/d"},
			{dir: "aa/b/c", want: "d"},
		}
		for _, tt := range tests {
			t.Run(tt.dir, func(t *testing.T) {
				s, err := fs.Sub(p, tt.dir)
				require.NoError(t, err)
				fstest.TestFS(s, tt.want)
			})
		}
	})

	t.Run("SubNotExist", func(t *testing.T) {
		tests := []struct {
			dir string
		}{
			{dir: "a"},
			{dir: "b"},
		}
		for _, tt := range tests {
			t.Run(tt.dir, func(t *testing.T) {
				_, err := fs.Sub(p, tt.dir)
				assert.ErrorAs(t, err, &fs.ErrNotExist)
			})
		}
	})

	t.Run("ReadDir", func(t *testing.T) {
		tests := []struct {
			dir  string
			want []string
		}{
			{dir: "aa", want: []string{"b"}},
			{dir: "aa/b", want: []string{"c"}},
			{dir: "aa/b/c", want: []string{"d"}},
		}
		for _, tt := range tests {
			t.Run(tt.dir, func(t *testing.T) {
				files, err := fs.ReadDir(p, tt.dir)
				require.NoError(t, err)
				var got []string
				for _, f := range files {
					got = append(got, f.Name())
				}
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("ReadDirNotExist", func(t *testing.T) {
		tests := []struct {
			dir string
		}{
			{dir: "a"},
			{dir: "b"},
		}
		for _, tt := range tests {
			t.Run(tt.dir, func(t *testing.T) {
				_, err := fs.ReadDir(p, tt.dir)
				assert.ErrorAs(t, err, &fs.ErrNotExist)
			})
		}
	})

	t.Run("Stat", func(t *testing.T) {
		tests := []struct {
			name  string
			isDir bool
		}{
			{name: "aa", isDir: true},
			{name: "aa/b", isDir: true},
			{name: "aa/b/c", isDir: true},
			{name: "aa/b/c/d", isDir: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := fs.Stat(p, tt.name)
				require.NoError(t, err)
				assert.Equal(t, tt.isDir, got.IsDir())
			})
		}
	})

	t.Run("StatNotExist", func(t *testing.T) {
		tests := []struct {
			name string
		}{
			{name: "a"},
			{name: "b"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := fs.Stat(p, tt.name)
				assert.ErrorAs(t, err, &fs.ErrNotExist)
			})
		}
	})
}
