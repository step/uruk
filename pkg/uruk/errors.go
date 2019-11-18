package uruk

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/step/saurontypes"
)

type ContainerCreationError struct {
	UrukMessage saurontypes.UrukMessage
	err         error
}

func (cce ContainerCreationError) Error() string {
	msg := cce.UrukMessage
	return fmt.Sprintf("Unable to create container %s for %s\n%s", msg.ImageName, msg.RepoLocation, cce.err.Error())
}

type CopyToContainerError struct {
	UrukMessage saurontypes.UrukMessage
	Source      string
	Destination string
	err         error
}

func (ctce CopyToContainerError) Error() string {
	msg := ctce.UrukMessage
	return fmt.Sprintf("Unable to copy from %s to %s:%s\n%s", ctce.Source, msg.ImageName, ctce.Destination, ctce.err.Error())
}

type StartContainerError struct {
	Response    container.ContainerCreateCreatedBody
	UrukMessage saurontypes.UrukMessage
	err         error
}

func (sce StartContainerError) Error() string {
	msg := sce.UrukMessage
	return fmt.Sprintf("Unable to start container of image %s with id %s for %s\nWarnings: %s\n%s",
		msg.ImageName,
		sce.Response.ID,
		msg.RepoLocation,
		sce.Response.Warnings,
		sce.err.Error())
}

type CopyFromContainerError struct {
	UrukMessage saurontypes.UrukMessage
	src         string
	ContainerId string
	err         error
}

func (cfce CopyFromContainerError) Error() string {
	msg := cfce.UrukMessage
	return fmt.Sprintf("Unable to copy %s from container of image %s with id %s\n%s",
		cfce.src,
		msg.ImageName,
		cfce.ContainerId,
		cfce.err.Error())
}

func actualErrStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// GzipReaderCreateError is typically returned when gzip.NewReader() returns
// an error. The location is simply a string that helps you identify where
// the error occcurred. Usually a function name.
type GzipReaderCreateError struct {
	actualErr error
	location  string
}

// Error returns a string that reports the location and the original error
func (g GzipReaderCreateError) Error() string {
	return fmt.Sprintf("Unable to create gzip reader at %s\n%s", g.location, actualErrStr(g.actualErr))
}

// GzipReaderCloseError is typically returned when a gzip.Reader.Close() returns
// an error. The location is simply a string that helps you identify where
// the error occcurred. Usually a function name.
type GzipReaderCloseError struct {
	actualErr error
	location  string
}

// Error returns a string that reports the location and the original error
func (g GzipReaderCloseError) Error() string {
	return fmt.Sprintf("Unable to close gzip reader at %s\n%s", g.location, actualErrStr(g.actualErr))
}

// TarHeaderError is typically returned when tar.Header.Next() returns an error.
// The location is simply a string that helps you identify where
// the error occcurred. Usually a function name.
type TarHeaderError struct {
	actualErr error
	location  string
}

// Error returns a string that reports the location and the original error
func (t TarHeaderError) Error() string {
	return fmt.Sprintf("Unable to read tar header %s\n%s", t.location, actualErrStr(t.actualErr))
}

// ExtractionError is returned when Extractor.ExtractFile returns an error
// Typically this error wraps one or more errors that might have occurred
// while either untarring the file or writing it. It accepts a file name
// the permission mode, the original error returned and a location. The
// location is simply a string that helps you identify where the error
// occurreed. Usually a function name.
type ExtractionError struct {
	name      string
	mode      int64
	actualErr error
	location  string
}

// Error returns a string that reports the name from the tar header, the mode, the
// original error returned and a location. The location is simply a string that helps you
// identify where the error occurred. Usually a function name.
func (e ExtractionError) Error() string {
	return fmt.Sprintf("Unable to extract %s [%o] at %s\n%s", e.name, e.mode, e.location, actualErrStr(e.actualErr))
}

// FileOpenError is typically returned when os.OpenFile returns an error
// It accepts a file name that includes the absolute path, the file name
// without the leading path, the permission mode, the original error returned
// and a location. The location is simply a string that helps you identify where
// the error occurreed. Usually a function name.
type FileOpenError struct {
	fileName  string
	name      string
	mode      int64
	actualErr error
	location  string
}

// Error returns a string that reports the absolute name of the file, the relative name ,
// the mode, the original error returned and a location. The location is simply a string
// that helps you identify where the error occurred. Usually a function name.
func (f FileOpenError) Error() string {
	return fmt.Sprintf("Unable to open %s(%s) [%o] at %s\n%s",
		f.fileName, f.name, f.mode, f.location, actualErrStr(f.actualErr))
}

// FileCopyError is typically returned when os.Copy returns an error
// The error needs to be named something better however since it is an error
// that relates to io.Writer and io.Reader and has nothing to do with ifiles
// It accepts a src, a dest being copied to, the number of bytes written as a part
// of copy and the error so far, the original error returned and a location.
// The location is simply a string that helps you identify where
// the error occurreed. Usually a function name.
type FileCopyError struct {
	src         string
	dest        string
	bytesCopied int64
	actualErr   error
	location    string
}

// Error returns a string that reports the src and destination being copied between
// the number of bytes that were copied, the original error returned and a location.
// The location is simply a string that helps you identify where the error occurred.
// Usually a function name.
func (f FileCopyError) Error() string {
	return fmt.Sprintf("Unable to copy from %s to %s\nCopied %d bytes at %s\n%s",
		f.src, f.dest, f.bytesCopied, f.location, actualErrStr(f.actualErr))
}

// FileCloseError is typically returned when file.Close() returns an error.
// Since file.Close() is often called in a defer block, two errors are provided
// actualErrr is the error that occurred while closing the file and pastErrors
// are errors that might have happened while opening a file for instance.
// The location is simply a string that helps you identify where the error occurred.
type FileCloseError struct {
	fileName   string
	actualErr  error
	pastErrors error
	location   string
}

// Error returns a string that reports the file that returned an error while
// closing. If pastErrors were included it reports that as well along with
// the location
func (c FileCloseError) Error() string {
	pastErrorsStr := ""
	if c.pastErrors != nil {
		pastErrorsStr = c.pastErrors.Error()
	}
	return fmt.Sprintf("Unable to close %s at %s\n%s\n%s", c.fileName, c.location, actualErrStr(c.actualErr), pastErrorsStr)
}

// MakeDirError is typically returned when os.MakeDir is called. The directory
// name and actual error which os.MakeDir returns are expected. The location
// is simply a string that helps you identify where the error occurred.
// Usually a function name.
type MakeDirError struct {
	dirName   string
	actualErr error
	location  string
}

// Error returns a string that reports which directory reported an error while
// creating it along with the actual error returned from os.MakeDir
func (d MakeDirError) Error() string {
	return fmt.Sprintf("Unable to make directory %s at %s\n%s", d.dirName, d.location, actualErrStr(d.actualErr))
}
