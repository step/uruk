package tarutils

import (
	"archive/tar"
	"os"
	"path/filepath"
	"strings"
)

type Tarable interface {
	ShouldTar(src, fileName string, fi os.FileInfo) bool
	GetHeader(src, fileName string, fi os.FileInfo) (*tar.Header, error)
}

type DefaultTarable struct {
	prefixPath string
}

func (d DefaultTarable) ShouldTar(src string, fileName string, fi os.FileInfo) bool {
	return true
}

func (d DefaultTarable) GetHeader(src string, fileName string, fi os.FileInfo) (*tar.Header, error) {
	baseStripped := strings.Replace(fileName, src, "", -1)

	header, err := tar.FileInfoHeader(fi, fi.Name())

	if err != nil {
		return nil, err
	}

	newFilename := filepath.Join(d.prefixPath, baseStripped)
	header.Name = newFilename
	return header, nil
}

func NewDefaultTarable(prefixPath string) Tarable {
	return DefaultTarable{
		prefixPath: prefixPath,
	}
}
