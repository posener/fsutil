package fsutil

import (
	"errors"
	"io/fs"
)

// UnionFS implements union of multiple filesystems. When used, the first filesystem that has the
// requested content is the one used. For example, if the same file exists in the first and second
// given filesystem, the file in the second will be hidden by the file in the first.
type UnionFS []fs.FS

// Open returns the file from the first filesystem that does not return an ErrNotExist. If the file
// does not exist in any filesystem, an ErrNotExist will be returned.
func (u UnionFS) Open(name string) (fs.File, error) {
	err := Check(name)
	if err != nil {
		return nil, err
	}
	for _, i := range u {
		f, err := i.Open(name)
		if !errors.Is(err, fs.ErrNotExist) {
			return f, err
		}
	}
	return nil, fs.ErrNotExist
}

// ReadDir is implemented in case that any of the given filesystems has an optimized ReadDir
// implementation.
func (u UnionFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var list []fs.DirEntry
	for _, i := range u {
		l, err := fs.ReadDir(i, name)
		if err != nil {
			return nil, err
		}
		list = append(list, l...)
	}
	return list, nil
}

// ReadFile is implemented in case that any of the given filesystems has an optimized ReadFile
// implementation.
func (u UnionFS) ReadFile(name string) ([]byte, error) {
	for _, i := range u {
		f, err := fs.ReadFile(i, name)
		if !errors.Is(err, fs.ErrNotExist) {
			return f, err
		}
	}
	return nil, fs.ErrNotExist
}

// Stat is implemented in case that any of the given filesystems has an optimized Stat
// implementation.
func (u UnionFS) Stat(name string) (fs.FileInfo, error) {
	for _, i := range u {
		s, err := fs.Stat(i, name)
		if !errors.Is(err, fs.ErrNotExist) {
			return s, err
		}
	}
	return nil, fs.ErrNotExist
}

// Sub is implemented in case that any of the given filesystems has an optimized Sub
// implementation.
func (u UnionFS) Sub(dir string) (fs.FS, error) {
	var sub UnionFS
	for _, i := range u {
		s, err := fs.Sub(i, dir)
		if err != nil {
			return nil, err
		}
		sub = append(sub, s)
	}
	return sub, nil
}

// Glob is implemented in case that any of the given filesystems has an optimized Glob
// implementation.
func (u UnionFS) Glob(pattern string) ([]string, error) {
	var names []string
	for _, i := range u {
		s, err := fs.Glob(i, pattern)
		if err != nil {
			return nil, err
		}
		names = append(names, s...)
	}
	return names, nil
}
