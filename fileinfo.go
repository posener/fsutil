package fsutil

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	VName    string      // base name of the file
	VSize    int64       // length in bytes for regular files; system-dependent for others
	VMode    fs.FileMode // file mode bits
	VModTime time.Time   // modification time
	VIsDir   bool        // abbreviation for Mode().IsDir()
	VSys     interface{} // underlying data source (can return nil)
}

func (f FileInfo) Name() string       { return f.VName }
func (f FileInfo) Size() int64        { return f.VSize }
func (f FileInfo) Mode() fs.FileMode  { return f.VMode }
func (f FileInfo) ModTime() time.Time { return f.VModTime }
func (f FileInfo) IsDir() bool        { return f.VIsDir }
func (f FileInfo) Sys() interface{}   { return f.VSys }

type Dir string

func (d Dir) Name() string { return string(d) }

func (d Dir) IsDir() bool { return true }

func (d Dir) Type() fs.FileMode { return fs.ModeDir }

func (d Dir) Info() (fs.FileInfo, error) {
	return FileInfo{VName: string(d), VMode: fs.ModeDir, VIsDir: true}, nil
}

func (d Dir) Stat() (fs.FileInfo, error) {
	return FileInfo{VName: string(d), VMode: fs.ModeDir, VIsDir: true}, nil
}
func (d Dir) Read([]byte) (int, error) { return 0, nil }
func (d Dir) Close() error             { return nil }
