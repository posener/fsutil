package fsutil

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// PrefixFS adds a prefix to the filesystem.
type PrefixFS struct {
	f      fs.FS
	prefix pathFS
}

// AddPrefix adds a given prefix to the given filesystem. For example, a given prefix 'foo', and a
// filesystem with a file 'bar'. In the returned filesystem, the file will be available on
// 'foo/bar'.
func AddPrefix(f fs.FS, prefix string) PrefixFS {
	return PrefixFS{
		f:      f,
		prefix: pathFS(prefix),
	}
}

func (p PrefixFS) Open(name string) (fs.File, error) {
	err := Check(name)
	if err != nil {
		return nil, err
	}
	if len(name) <= len(p.prefix) {
		return p.prefix.Open(name)
	}
	name, err = trimPrefix(string(p.prefix), name)
	if err != nil {
		return nil, err
	}
	return p.f.Open(name)
}

// ReadDir is implemented in case that the given filesystem has an optimized ReadDir
// implementation.
func (p PrefixFS) ReadDir(name string) ([]fs.DirEntry, error) {
	err := Check(name)
	if err != nil {
		return nil, err
	}
	if len(name) < len(p.prefix) {
		return p.prefix.ReadDir(name)
	}
	name, err = trimPrefix(string(p.prefix), name)
	if err != nil {
		return nil, err
	}
	if name == "" {
		name = "."
	}
	return fs.ReadDir(p.f, name)
}

// ReadFile is implemented in case that the given filesystem has an optimized ReadFile
// implementation.
func (p PrefixFS) ReadFile(name string) ([]byte, error) {
	err := Check(name)
	if err != nil {
		return nil, err
	}
	name, err = trimPrefix(string(p.prefix), name)
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(p.f, name)
}

// Stat is implemented in case that the given filesystem has an optimized Stat implementation.
func (p PrefixFS) Stat(name string) (fs.FileInfo, error) {
	err := Check(name)
	if err != nil {
		return nil, err
	}
	if len(name) <= len(p.prefix) {
		return p.prefix.Stat(name)
	}
	name, err = trimPrefix(string(p.prefix), name)
	if err != nil {
		return nil, err
	}
	return fs.Stat(p.f, name)
}

// Sub is implemented in case that the given filesystem has an optimized Sub implementation.
func (p PrefixFS) Sub(dir string) (fs.FS, error) {
	err := Check(dir)
	if err != nil {
		return nil, err
	}
	if len(dir) < len(p.prefix) {
		sub, err := p.prefix.Sub(dir)
		if err != nil {
			return nil, err
		}
		p.prefix = sub.(pathFS)
		return p, nil
	}

	dir, err = trimPrefix(string(p.prefix), dir)
	if err != nil {
		return nil, err
	}
	if dir == "" {
		return p.f, nil
	}
	return fs.Sub(p.f, dir)
}

// pathFS is a helper filesystem for prefixFS. It represents the path of the prefix.
type pathFS string

func (p pathFS) Open(name string) (fs.File, error) {
	var err error
	name, err = trimPrefix(name, string(p))
	if err != nil {
		return nil, err
	}
	if name == "" {
		return Dir(filepath.Base(string(p))), nil
	}
	name = firstPart(name)
	return Dir(name), nil
}

func (p pathFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var err error
	name, err = trimPrefix(name, string(p))
	if err != nil {
		return nil, err
	}
	name = firstPart(name)
	return []fs.DirEntry{Dir(name)}, nil
}

func (p pathFS) Stat(name string) (fs.FileInfo, error) {
	var err error
	name, err = trimPrefix(name, string(p))
	if err != nil {
		return nil, err
	}
	if name == "" {
		return Dir(filepath.Base(string(p))).Info()
	}
	name = firstPart(name)
	return Dir(filepath.Base(name)).Info()
}

func (p pathFS) Sub(dir string) (fs.FS, error) {
	var err error
	sub, err := trimPrefix(dir, string(p))
	if err != nil {
		return nil, err
	}
	return pathFS(sub), nil
}

func trimPrefix(prefix string, name string) (string, error) {
	if !strings.HasPrefix(name, prefix) {
		return "", fs.ErrNotExist
	}
	name = name[len(prefix):]
	if len(name) == 0 {
		return "", nil
	}
	if name[0] != '/' {
		return "", fs.ErrNotExist
	}
	return name[1:], nil
}

func firstPart(path string) string {
	if i := strings.IndexByte(path, '/'); i > 0 {
		path = path[:i]
	}
	return path
}
