package frontend

import (
	"net/http"
	"strings"

	challenger "github.com/P0H1ng/fronted_player/dist"
	manager "github.com/P0H1ng/fronted_manager/dist"
)

type frontendFS struct {
	frontendFS http.FileSystem
	managerFS  http.FileSystem
}

// FS is the filesystem of the frontend.
func FS() *frontendFS {
	return &frontendFS{
		frontendFS: challenger.New(),
		managerFS:  manager.New(),
	}
}

// Open: open file.
func (f *frontendFS) Open(name string) (http.File, error) {
	if strings.HasPrefix(name, "/manager") {
		return f.managerFS.Open(name)
	}
	return f.frontendFS.Open(name)
}

// Exists: check if the file exist.
func (f *frontendFS) Exists(prefix string, filePath string) bool {
	if strings.HasPrefix(filePath, "/manager") {
		if _, err := f.managerFS.Open(filePath); err != nil {
			return false
		}
		return true
	}
	if _, err := f.frontendFS.Open(filePath); err != nil {
		return false
	}
	return true
}
