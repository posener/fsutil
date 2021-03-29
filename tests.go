package fsutil

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertContent(t *testing.T, f fs.FS, path string, wantContent string) {
	t.Helper()
	got, err := fs.ReadFile(f, path)
	require.NoError(t, err)
	assert.Equal(t, wantContent, string(got))
}

func assertNotExists(t *testing.T, f fs.FS, path string) {
	t.Helper()
	_, err := f.Open(path)
	assert.ErrorIs(t, err, fs.ErrNotExist)

	_, err = fs.Stat(f, path)
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
