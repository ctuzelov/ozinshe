package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func ProcessSaving(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func DeleteDirectory(DirectoryPath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current dir: %v", err)
	}

	DirectoryPath = currentDir + "/internal/" + DirectoryPath

	err = os.RemoveAll(DirectoryPath)
	if err != nil {
		return fmt.Errorf("error deleting directory '%s': %v", DirectoryPath, err)
	}
	return nil
}
