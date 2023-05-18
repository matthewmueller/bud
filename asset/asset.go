package asset

import "io/fs"

func FS(fsys fs.FS) fs.FS {
	return fsys
}
