package kasset

import (
	"os"
	"path/filepath"
)

type AssetFS interface {
	Save(string, []byte) error
	Read(string) ([]byte, error)
	Delete(string) error
}

type simpleFS struct {
	BasicPath string
}

func NewSimpleFS(basicPath string) *simpleFS {
	s := new(simpleFS)
	s.BasicPath = basicPath
	return s
}

func (sfs *simpleFS) Save(name string, bs []byte) error {
	fullPath := filepath.Join(sfs.BasicPath, name)
	return os.WriteFile(fullPath, bs, 0644)
}

func (sfs *simpleFS) Read(name string) ([]byte, error) {
	fullPath := filepath.Join(sfs.BasicPath, name)
	bs, e := os.ReadFile(fullPath)
	return bs, e
}

func (sfs *simpleFS) Delete(name string) error {
	fullPath := filepath.Join(sfs.BasicPath, name)
	return os.Remove(fullPath)
}
