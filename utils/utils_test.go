package utils

import (
	"os"
	"path"
	"testing"
)

func TestCopyDir(T *testing.T) {
	// Test cases for CopyDir function

	// Create a temp directory for test
	tempSourceDir, err := os.MkdirTemp("", "test")
	if err != nil {
		T.Fatal(err)
	}
	tempTargetDir, err := os.MkdirTemp("", "test")
	if err != nil {
		T.Fatal(err)
	}
	// defer func() { // Clean up
	// 	if err := os.RemoveAll(tempSourceDir); err != nil {
	// 		T.Fatal(err)
	// 	}
	// 	if err := os.RemoveAll(tempTargetDir); err != nil {
	// 		T.Fatal(err)
	// 	}
	// }()
	// Create contents
	var file_paths_and_contents = []struct {
		path string
		data string
	}{
		{"file_Lv0", "content_Lv0"},
		{"file_Lv1", "content_Lv1"},
		{"file_Lv2", "content_Lv2"},
		{"file_Lv3", "content_Lv3"},
	}

	// create dirs
	if err := os.MkdirAll(path.Join(tempSourceDir, "folder_Lv1", "folder_Lv2", "Folder_Lv3", "Folder_Lv4"), 0700); err != nil {
		T.Fatal(err)
	}
	// create files
	for _, file := range file_paths_and_contents {
		file_path := path.Join(tempSourceDir, file.path)
		if err := os.WriteFile(file_path, []byte(file.data), 0700); err != nil {
			T.Fatal(err)
		}
	}

	// Copy the directory to the temp directory
	err = CopyDir(tempSourceDir, tempTargetDir)
	if err != nil {
		T.Fatal(err)
	}

	// TODO: Check if the contents of the source directory are successfully copied.
}
