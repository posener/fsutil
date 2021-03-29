package fsutil

import "io/fs"

// Check checks if the name is `fs.ValidPath`, and returns the appropriate error accordingly.
func Check(name string) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	return nil
}
