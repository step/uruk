package uruk

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor is an interface that can be used whenever one
// wants to extract a file or a directory based on a tar header
// It is left to the caller that calls the interface to decide
// if the header is that of a file or a directory
type Extractor interface {
	ExtractFile(tar.Header, io.Reader) error
	ExtractDir(tar.Header, io.Reader) error
	GetBasePath() string
	String() string
}

// ExtractorGenerator is an interface that is used to generate
// an extractor
type ExtractorGenerator interface {
	Generate(...string) Extractor
	String() string
}

// DefaultExtractor is a struct implementing the Extractor interface.
// It has a src field that indicates the base path where one wants to
// extract files and directories to.
type DefaultExtractor struct {
	src string
}

func stripLeadingComponent(path string) string {
	pathComponents := strings.Split(path, string(filepath.Separator))
	return filepath.Join(pathComponents[1:]...)
}

// ExtractFile extracts the given file under src specified
// in DefaultExtractor
func (extractor DefaultExtractor) ExtractFile(header tar.Header, reader io.Reader) (rerr error) {
	// Open file and defer file.Close()
	// TODO: create a more general file ignore mechanism
	if header.Name == "pax_global_header" {
		return nil
	}

	location := "DefaultExtractor.ExtractFile"
	actualPath := stripLeadingComponent(header.Name)
	fileName := filepath.Join(extractor.src, actualPath)
	file, ferr := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())

	// Handle defer in an anonymous func
	defer func() {
		err := file.Close()
		if err != nil {
			rerr = FileCloseError{fileName, err, rerr, location}
		}
	}()

	// Unable to open file
	if ferr != nil {
		return FileOpenError{fileName, header.Name, header.Mode, ferr, location}
	}

	// Copy file
	numBytesCopied, err := io.Copy(file, reader)

	// Unable to copy file
	if err != nil {
		return FileCopyError{header.Name, fileName, numBytesCopied, err, location}
	}

	return nil
}

// ExtractDir extracts the given dir under src specified
// in DefaultExtractor
func (extractor DefaultExtractor) ExtractDir(header tar.Header, reader io.Reader) error {
	// Create directory in src
	actualPath := stripLeadingComponent(header.Name)
	dirName := filepath.Join(extractor.src, actualPath)
	err := os.MkdirAll(dirName, header.FileInfo().Mode())
	if err != nil {
		return MakeDirError{dirName, err, "DefaultExtractor.ExtractDir"}
	}

	return nil
}

// GetBasePath is a method to get the base path from DefaultExtractor
func (extractor DefaultExtractor) GetBasePath() string {
	return extractor.src
}

func (extractor DefaultExtractor) String() string {
	return extractor.GetBasePath()
}

// NewDefaultExtractor creates an instance of DefaultExtractor with the specified src
func NewDefaultExtractor(src string) Extractor {
	return DefaultExtractor{src}
}

// DefaultExtractorGenerator is a simple implementation
// of ExtractorGenerator
type DefaultExtractorGenerator struct {
	Src string
}

// Generate is a method on DefaultExtractorGenerator
// to conveniently generate a new DefaultExtractorGenerator
func (d DefaultExtractorGenerator) Generate(args ...string) Extractor {
	relativePath := filepath.Join(args...)
	dir := filepath.Join(d.Src, relativePath)
	return NewDefaultExtractor(dir)
}

func (d DefaultExtractorGenerator) String() string {
	return fmt.Sprintf("DefaultExtractorGenerator: %s\n", d.Src)
}
