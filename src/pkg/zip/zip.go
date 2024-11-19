package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Unzip(dst, filename string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, f := range archive.File {
		filepath := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(filepath, f.Mode()); err != nil {
				return fmt.Errorf("could not create directory %s: %v", f.Name, err)
			}
			continue
		}

		if err = extractFile(f, filepath); err != nil {
			return fmt.Errorf("could not extract file %s: %v", f.Name, err)
		}
	}
	return nil
}

func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return err
	}

	err = os.Chmod(destPath, f.Mode())
	if err != nil {
		return err
	}

	return nil
}
