package tarutils

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Tar(src string, buffer io.Writer, tarable Tarable) {
	tarWriter := tar.NewWriter(buffer)
	defer tarWriter.Close()

	filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !tarable.ShouldTar(src, fileName, fi) {
			return nil
		}

		header, err := tarable.GetHeader(src, fileName, fi)
		if err != nil {
			return fmt.Errorf("Unable to get header for %s", fileName)
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("Unable to write header for %s", fi.Name())
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("Unable to open %s", fileName)
		}

		if _, err := io.Copy(tarWriter, file); err != nil {
			return fmt.Errorf("Unable to tar %s", fileName)
		}

		fmt.Println(" +", header.Name)
		if err := file.Close(); err != nil {
			return fmt.Errorf("Unable to close %s", fileName)
		}

		return nil
	})
}
