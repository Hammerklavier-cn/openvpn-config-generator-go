package utils

import (
	"fmt"
	"io"
	"os"
	"path"
)

func CopyDir(sourceDir, targetDir string) error {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("Failed to read easy-rsa source directory: %w", err)
	}
	// fmt.Println("source dir:", sourceDir, "target dir:", targetDir)
	// fmt.Println("Files and Dirs:", files)

	for _, file := range files {
		sourcePath := path.Join(sourceDir, file.Name())
		targetPath := path.Join(targetDir, file.Name())

		if file.IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("Failed to create target directory: %w", err)
			}
			if err := CopyDir(sourcePath, targetPath); err != nil {
				return err
			}
		} else {
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				return fmt.Errorf("Failed to open source file: %w", err)
			}
			defer sourceFile.Close()

			targetFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("Failed to create target file: %w", err)
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, sourceFile); err != nil {
				return fmt.Errorf("Failed to copy file contents: %w", err)
			}
		}
	}
	return nil
}
